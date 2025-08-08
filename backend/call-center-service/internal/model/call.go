package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CallStatus 呼叫状态
type CallStatus string

const (
	CallStatusPending    CallStatus = "pending"    // 待呼叫
	CallStatusDialing    CallStatus = "dialing"    // 拨号中
	CallStatusRinging    CallStatus = "ringing"    // 振铃中
	CallStatusConnected  CallStatus = "connected"  // 已接通
	CallStatusHangup     CallStatus = "hangup"     // 已挂断
	CallStatusFailed     CallStatus = "failed"     // 呼叫失败
	CallStatusNoAnswer   CallStatus = "no_answer"  // 无人接听
	CallStatusBusy       CallStatus = "busy"       // 忙线
	CallStatusCancelled  CallStatus = "cancelled"  // 已取消
)

// CallDirection 呼叫方向
type CallDirection string

const (
	CallDirectionInbound  CallDirection = "inbound"  // 呼入
	CallDirectionOutbound CallDirection = "outbound" // 呼出
)

// DialMode 拨号模式
type DialMode string

const (
	DialModePredictive DialMode = "predictive" // 预测式
	DialModePreview    DialMode = "preview"    // 预览式
	DialModeAuto       DialMode = "auto"       // 自动
	DialModeManual     DialMode = "manual"     // 手动
)

// Call 呼叫记录
type Call struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID     uuid.UUID      `gorm:"type:uuid;index" json:"project_id"`
	CampaignID    uuid.UUID      `gorm:"type:uuid;index" json:"campaign_id"`
	AgentID       *uuid.UUID     `gorm:"type:uuid;index" json:"agent_id"`
	CustomerPhone string         `gorm:"type:varchar(20);index" json:"customer_phone"`
	CallerID      string         `gorm:"type:varchar(20)" json:"caller_id"`
	Direction     CallDirection  `gorm:"type:varchar(20)" json:"direction"`
	Status        CallStatus     `gorm:"type:varchar(20);index" json:"status"`
	DialMode      DialMode       `gorm:"type:varchar(20)" json:"dial_mode"`
	StartTime     *time.Time     `json:"start_time"`
	ConnectTime   *time.Time     `json:"connect_time"`
	EndTime       *time.Time     `json:"end_time"`
	Duration      int            `json:"duration"`       // 通话时长（秒）
	WaitTime      int            `json:"wait_time"`      // 等待时长（秒）
	RecordingURL  string         `json:"recording_url"`
	CallSID       string         `gorm:"type:varchar(100);uniqueIndex" json:"call_sid"`
	Metadata      map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (c *Call) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// CallSession 呼叫会话
type CallSession struct {
	ID           uuid.UUID    `gorm:"type:uuid;primary_key" json:"id"`
	CallID       uuid.UUID    `gorm:"type:uuid;index" json:"call_id"`
	AgentID      uuid.UUID    `gorm:"type:uuid;index" json:"agent_id"`
	SessionType  string       `gorm:"type:varchar(20)" json:"session_type"` // agent, customer
	JoinTime     time.Time    `json:"join_time"`
	LeaveTime    *time.Time   `json:"leave_time"`
	State        string       `gorm:"type:varchar(20)" json:"state"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (cs *CallSession) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	return nil
}

// Recording 录音记录
type Recording struct {
	ID            uuid.UUID    `gorm:"type:uuid;primary_key" json:"id"`
	CallID        uuid.UUID    `gorm:"type:uuid;index" json:"call_id"`
	URL           string       `json:"url"`
	Duration      int          `json:"duration"`
	Size          int64        `json:"size"`
	Format        string       `gorm:"type:varchar(10)" json:"format"`
	StoragePath   string       `json:"storage_path"`
	TranscriptURL string       `json:"transcript_url"`
	CreatedAt     time.Time    `json:"created_at"`
}

func (r *Recording) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
