package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/cloudcallcenter/call-center-service/internal/service"
)

// AuditOptions 审计中间件选项

type AuditOptions struct {
	// SkipPaths 不记录的路径前缀
	SkipPaths []string
	// ActorResolver 从请求上下文提取 actor 信息（可选）
	ActorResolver func(c *gin.Context) (actorID *uuid.UUID, actorType string)
	// TraceIDResolver 从请求上下文提取 TraceID（可选）
	TraceIDResolver func(c *gin.Context) string
	// ResourceResolver 根据请求解析资源信息（可选）
	ResourceResolver func(c *gin.Context, status int) (resourceType string, resourceID *uuid.UUID, action string, level string)
}

func AuditLogger(svc service.AuditService, opts *AuditOptions) gin.HandlerFunc {
	if opts == nil {
		opts = &AuditOptions{}
	}
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, p := range opts.SkipPaths {
			if len(path) >= len(p) && path[:len(p)] == p {
				c.Next()
				return
			}
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		status := c.Writer.Status()
		method := c.Request.Method
		ua := c.Request.UserAgent()
		ip := c.ClientIP()

		var actorID *uuid.UUID
		var actorType string
		if opts.ActorResolver != nil {
			actorID, actorType = opts.ActorResolver(c)
		}

		var traceID string
		if opts.TraceIDResolver != nil {
			traceID = opts.TraceIDResolver(c)
		}

		var resourceType string
		var resourceID *uuid.UUID
		var action string
		level := "info"
		if opts.ResourceResolver != nil {
			resourceType, resourceID, action, level = opts.ResourceResolver(c, status)
		} else {
			// 基于HTTP方法的默认推断
			switch method {
			case "POST":
				action = "create"
			case "PUT", "PATCH":
				action = "update"
			case "DELETE":
				action = "delete"
			default:
				action = "get"
			}
		}

		entry := &service.AuditLogInput{
			ActorID:      actorID,
			ActorType:    actorType,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Level:        level,
			Method:       method,
			Path:         path,
			StatusCode:   status,
			IP:           ip,
			UserAgent:    ua,
			LatencyMs:    latency.Milliseconds(),
			TraceID:      traceID,
		}

		// 异步写入，避免阻塞请求
		go func() { _ = svc.Log(c.Request.Context(), entry) }()
	}
}