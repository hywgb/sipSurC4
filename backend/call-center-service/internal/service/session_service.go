package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	
	"github.com/cloudcallcenter/call-center-service/internal/model"
	"github.com/cloudcallcenter/call-center-service/internal/repository"
)

// sessionService 会话服务实现
type sessionService struct {
	sessionRepo repository.SessionRepository
	logger      *logrus.Logger
}

// NewSessionService 创建会话服务
func NewSessionService(sessionRepo repository.SessionRepository) SessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
		logger:      logrus.New(),
	}
}

// CreateSession 创建会话
func (s *sessionService) CreateSession(ctx context.Context, session *model.CallSession) error {
	session.JoinTime = time.Now()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	
	s.logger.Infof("Creating session for call %s, agent %s", session.CallID, session.AgentID)
	return s.sessionRepo.Create(ctx, session)
}

// GetSession 获取会话信息
func (s *sessionService) GetSession(ctx context.Context, sessionID uuid.UUID) (*model.CallSession, error) {
	return s.sessionRepo.GetByID(ctx, sessionID)
}

// GetActiveSessions 获取活跃会话
func (s *sessionService) GetActiveSessions(ctx context.Context) ([]*model.CallSession, error) {
	return s.sessionRepo.GetActiveSessions(ctx)
}

// GetSessionsByCall 根据呼叫ID获取会话
func (s *sessionService) GetSessionsByCall(ctx context.Context, callID uuid.UUID) ([]*model.CallSession, error) {
	return s.sessionRepo.GetByCallID(ctx, callID)
}

// EndSession 结束会话
func (s *sessionService) EndSession(ctx context.Context, sessionID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"leave_time":  now,
		"state":       "ended",
		"updated_at":  now,
	}
	
	s.logger.Infof("Ending session %s", sessionID)
	return s.sessionRepo.Update(ctx, sessionID, updates)
}
