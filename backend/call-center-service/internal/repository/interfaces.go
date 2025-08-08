package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// CallRepository 呼叫仓储接口
type CallRepository interface {
	Create(ctx context.Context, call *model.Call) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Call, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*model.Call, error)
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// AgentRepository 座席仓储接口
type AgentRepository interface {
	Create(ctx context.Context, agent *model.Agent) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Agent, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Agent, error)
	GetByWorkID(ctx context.Context, workID string) (*model.Agent, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*model.Agent, error)
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.AgentStatus) error
}

// SessionRepository 会话仓储接口
type SessionRepository interface {
	Create(ctx context.Context, session *model.CallSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.CallSession, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*model.CallSession, error)
	GetByCallID(ctx context.Context, callID uuid.UUID) ([]*model.CallSession, error)
	GetActiveSessions(ctx context.Context) ([]*model.CallSession, error)
}

// RecordingRepository 录音仓储接口
type RecordingRepository interface {
	Create(ctx context.Context, recording *model.Recording) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Recording, error)
	GetByCallID(ctx context.Context, callID uuid.UUID) (*model.Recording, error)
	List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*model.Recording, error)
}

// AgentStatsRepository 座席统计仓储接口
type AgentStatsRepository interface {
	Create(ctx context.Context, stats *model.AgentStats) error
	GetByAgentIDAndDate(ctx context.Context, agentID uuid.UUID, date string) (*model.AgentStats, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	IncrementStats(ctx context.Context, agentID uuid.UUID, date string, field string, value int) error
}

// AgentScheduleRepository 座席排班仓储接口
type AgentScheduleRepository interface {
	Create(ctx context.Context, schedule *model.AgentSchedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.AgentSchedule, error)
	GetByAgentIDAndDateRange(ctx context.Context, agentID uuid.UUID, startDate, endDate string) ([]*model.AgentSchedule, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
