package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cloudcallcenter/call-center-service/internal/handler"
	"github.com/cloudcallcenter/call-center-service/internal/model"
	"github.com/cloudcallcenter/call-center-service/internal/repository"
	"github.com/cloudcallcenter/call-center-service/internal/service"
	"github.com/cloudcallcenter/call-center-service/pkg/dialer"
	"github.com/cloudcallcenter/call-center-service/pkg/router"
	"github.com/cloudcallcenter/call-center-service/pkg/middleware"
)

var (
	log = logrus.New()
)

func init() {
	// 设置日志格式
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	// 加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	
	viper.SetDefault("server.http.port", 8080)
	viper.SetDefault("server.grpc.port", 9090)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("redis.addr", "localhost:6379")

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Config file not found, using defaults: %v", err)
	}
}

func main() {
	// 初始化数据库
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移数据库结构
	if err := db.AutoMigrate(
		&model.Call{},
		&model.Agent{},
		&model.CallSession{},
		&model.Recording{},
		&model.AuditLog{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化仓储层
	callRepo := repository.NewCallRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	recordingRepo := repository.NewRecordingRepository(db)

	// 初始化拨号引擎
	dialerEngine := dialer.NewDialer()
	
	// 初始化路由引擎
	routerEngine := router.NewRouter()

	// 初始化服务层
	callService := service.NewCallService(callRepo, dialerEngine, routerEngine, recordingRepo)
	agentService := service.NewAgentService(agentRepo)
	sessionService := service.NewSessionService(sessionRepo)
	auditRepo := repository.NewAuditLogRepository(db)
	auditService := service.NewAuditService(auditRepo)

	// 初始化HTTP服务器
	httpServer := setupHTTPServer(callService, agentService, sessionService, auditService)
	
	// 初始化gRPC服务器
	grpcServer := setupGRPCServer(callService, agentService)

	// 启动HTTP服务器
	go func() {
		httpPort := viper.GetInt("server.http.port")
		log.Infof("Starting HTTP server on port %d", httpPort)
		if err := httpServer.Run(fmt.Sprintf(":%d", httpPort)); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// 启动gRPC服务器
	go func() {
		grpcPort := viper.GetInt("server.grpc.port")
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		if err != nil {
			log.Fatalf("Failed to listen on port %d: %v", grpcPort, err)
		}
		log.Infof("Starting gRPC server on port %d", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// 启动Prometheus metrics服务器
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		metricsPort := viper.GetInt("server.metrics.port")
		if metricsPort == 0 {
			metricsPort = 9091
		}
		log.Infof("Starting metrics server on port %d", metricsPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil); err != nil {
			log.Errorf("Failed to start metrics server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// 优雅关闭
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
 
	grpcServer.GracefulStop()
	
	log.Info("Server exited")
}

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		viper.GetString("database.host"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
		viper.GetInt("database.port"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func setupHTTPServer(
	callService service.CallService,
	agentService service.AgentService,
	sessionService service.SessionService,
	auditService service.AuditService,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	
	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 审计日志中间件（跳过健康检查和指标）
	r.Use(middleware.AuditLogger(auditService, &middleware.AuditOptions{SkipPaths: []string{"/health", "/metrics"}}))
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API路由
	api := r.Group("/api/v1")
	{
		// 呼叫相关API
		callHandler := handler.NewCallHandler(callService)
		api.POST("/calls", callHandler.CreateCall)
		api.GET("/calls", callHandler.ListCalls)
		api.GET("/calls/:id", callHandler.GetCall)
		api.PUT("/calls/:id", callHandler.UpdateCall)
		api.POST("/calls/:id/hangup", callHandler.HangupCall)
		api.GET("/calls/:id/recording", callHandler.GetRecording)

		// 座席相关API
		agentHandler := handler.NewAgentHandler(agentService)
		api.GET("/agents", agentHandler.ListAgents)
		api.GET("/agents/:id", agentHandler.GetAgent)
		api.PUT("/agents/:id/status", agentHandler.UpdateAgentStatus)
		api.GET("/agents/:id/stats", agentHandler.GetAgentStats)

		// 会话相关API
		sessionHandler := handler.NewSessionHandler(sessionService)
		api.GET("/sessions", sessionHandler.ListSessions)
		api.GET("/sessions/:id", sessionHandler.GetSession)
		api.GET("/sessions/active", sessionHandler.GetActiveSessions)

		// 审计日志API
		auditHandler := handler.NewAuditHandler(auditService)
		api.GET("/audits", auditHandler.ListAudits)
	}

	return r
}

func setupGRPCServer(
	callService service.CallService,
	agentService service.AgentService,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 10), // 10MB
	}
	
	server := grpc.NewServer(opts...)
	
	// 注册gRPC服务
	// pb.RegisterCallServiceServer(server, callService)
	// pb.RegisterAgentServiceServer(server, agentService)
	
	return server
}
