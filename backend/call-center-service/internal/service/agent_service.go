package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	
	"github.com/cloudcallcenter/call-center-service/internal/model"
	"github.com/cloudcallcenter/call-center-service/internal/repository"
)

// agentService 座席服务实现
type agentService struct {
	agentRepo    repository.AgentRepository
	statsRepo    repository.AgentStatsRepository
	scheduleRepo repository.AgentScheduleRepository
	logger      *logrus.Logger
}

// NewAgentService 创建座席服务
func NewAgentService(
	agentRepo repository.AgentRepository,
	statsRepo repository.AgentStatsRepository,
	scheduleRepo repository.AgentScheduleRepository,
) AgentService {
	return &agentService{
		agentRepo:    agentRepo,
		statsRepo:    statsRepo,
		scheduleRepo: scheduleRepo,
		logger:      logrus.New(),
	}
}

// CreateAgent 创建座席
func (s *agentService) CreateAgent(ctx context.Context, agent *model.Agent) error {
	agent.Status = model.AgentStatusOffline
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()
	return s.agentRepo.Create(ctx, agent)
}

// GetAgent 获取座席信息
func (s *agentService) GetAgent(ctx context.Context, agentID uuid.UUID) (*model.Agent, error) {
	return s.agentRepo.GetByID(ctx, agentID)
}

// UpdateAgent 更新座席信息
func (s *agentService) UpdateAgent(ctx context.Context, agentID uuid.UUID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return s.agentRepo.Update(ctx, agentID, updates)
}

// ListAgents 查询座席列表
func (s *agentService) ListAgents(ctx context.Context, filter *AgentFilter, page, pageSize int) ([]*model.Agent, int64, error) {
	query := make(map[string]interface{})
	
	if filter.Status != nil {
		query["status"] = *filter.Status
	}
	if filter.Level != nil {
		query["level"] = *filter.Level
	}
	
	offset := (page - 1) * pageSize
	
	agents, err := s.agentRepo.List(ctx, query, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := s.agentRepo.Count(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	return agents, total, nil
}

// UpdateAgentStatus 更新座席状态
func (s *agentService) UpdateAgentStatus(ctx context.Context, agentID uuid.UUID, status model.AgentStatus) error {
	return s.agentRepo.UpdateStatus(ctx, agentID, status)
}

// GetAgentStats 获取座席统计
func (s *agentService) GetAgentStats(ctx context.Context, agentID uuid.UUID, date string) (*model.AgentStats, error) {
	if s.statsRepo == nil {
		// 返回基本占位数据
		return &model.AgentStats{AgentID: agentID, Date: date}, nil
	}
	stats, err := s.statsRepo.GetByAgentIDAndDate(ctx, agentID, date)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// GetAgentSchedule 获取座席排班
func (s *agentService) GetAgentSchedule(ctx context.Context, agentID uuid.UUID, startDate, endDate string) ([]*model.AgentSchedule, error) {
	if s.scheduleRepo == nil {
		return []*model.AgentSchedule{}, nil
	}
	return s.scheduleRepo.GetByAgentIDAndDateRange(ctx, agentID, startDate, endDate)
}

// Login 座席登录
func (s *agentService) Login(ctx context.Context, agentID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":        model.AgentStatusAvailable,
		"last_login_at": now,
		"updated_at":    now,
	}
	
	s.logger.Infof("Agent %s logged in", agentID)
	return s.agentRepo.Update(ctx, agentID, updates)
}

// Logout 座席登出
func (s *agentService) Logout(ctx context.Context, agentID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":         model.AgentStatusOffline,
		"last_logout_at": now,
		"updated_at":     now,
	}
	
	s.logger.Infof("Agent %s logged out", agentID)
	return s.agentRepo.Update(ctx, agentID, updates)
}
