package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// CallService 呼叫服务接口
type CallService interface {
	CreateCall(ctx context.Context, req *CreateCallRequest) (*model.Call, error)
	GetCall(ctx context.Context, callID uuid.UUID) (*model.Call, error)
	UpdateCall(ctx context.Context, callID uuid.UUID, updates map[string]interface{}) error
	ListCalls(ctx context.Context, filter *CallFilter, page, pageSize int) ([]*model.Call, int64, error)
	HangupCall(ctx context.Context, callID uuid.UUID) error
	GetCallRecording(ctx context.Context, callID uuid.UUID) (*model.Recording, error)
	GetCallStats(ctx context.Context, timeRange TimeRange) (*CallStats, error)
}

// AgentService 座席服务接口
type AgentService interface {
	CreateAgent(ctx context.Context, agent *model.Agent) error
	GetAgent(ctx context.Context, agentID uuid.UUID) (*model.Agent, error)
	UpdateAgent(ctx context.Context, agentID uuid.UUID, updates map[string]interface{}) error
	ListAgents(ctx context.Context, filter *AgentFilter, page, pageSize int) ([]*model.Agent, int64, error)
	UpdateAgentStatus(ctx context.Context, agentID uuid.UUID, status model.AgentStatus) error
	GetAgentStats(ctx context.Context, agentID uuid.UUID, date string) (*model.AgentStats, error)
	GetAgentSchedule(ctx context.Context, agentID uuid.UUID, startDate, endDate string) ([]*model.AgentSchedule, error)
	Login(ctx context.Context, agentID uuid.UUID) error
	Logout(ctx context.Context, agentID uuid.UUID) error
}

// SessionService 会话服务接口
type SessionService interface {
	CreateSession(ctx context.Context, session *model.CallSession) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*model.CallSession, error)
	GetActiveSessions(ctx context.Context) ([]*model.CallSession, error)
	GetSessionsByCall(ctx context.Context, callID uuid.UUID) ([]*model.CallSession, error)
	EndSession(ctx context.Context, sessionID uuid.UUID) error
}

// CreateCallRequest 创建呼叫请求
type CreateCallRequest struct {
	ProjectID     uuid.UUID              `json:"project_id"`
	CampaignID    uuid.UUID              `json:"campaign_id"`
	CustomerPhone string                 `json:"customer_phone"`
	CallerID      string                 `json:"caller_id"`
	Direction     model.CallDirection    `json:"direction"`
	DialMode      model.DialMode         `json:"dial_mode"`
	Priority      int                    `json:"priority"`
	Skills        []string               `json:"skills"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// CallStatus 呼叫状态（用于避免循环依赖）
type CallStatus = model.CallStatus

// CallFilter 呼叫过滤条件
type CallFilter struct {
	ProjectID  *uuid.UUID
	CampaignID *uuid.UUID
	AgentID    *uuid.UUID
	Status     *model.CallStatus
	Direction  *model.CallDirection
	StartDate  *time.Time
	EndDate    *time.Time
	Phone      *string
}

// AgentFilter 座席过滤条件
type AgentFilter struct {
	Status    *model.AgentStatus
	Skills    []string
	Level     *int
	Available *bool
}

// TimeRange 时间范围
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// CallStats 呼叫统计
type CallStats struct {
	TotalCalls      int64   `json:"total_calls"`
	AnsweredCalls   int64   `json:"answered_calls"`
	MissedCalls     int64   `json:"missed_calls"`
	FailedCalls     int64   `json:"failed_calls"`
	AvgCallDuration float64 `json:"avg_call_duration"`
	AvgWaitTime     float64 `json:"avg_wait_time"`
	AnswerRate      float64 `json:"answer_rate"`
	AbandonRate     float64 `json:"abandon_rate"`
}
