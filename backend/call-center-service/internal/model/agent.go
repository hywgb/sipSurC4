package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgentStatus 座席状态
type AgentStatus string

const (
	AgentStatusOffline    AgentStatus = "offline"    // 离线
	AgentStatusAvailable  AgentStatus = "available"  // 空闲
	AgentStatusBusy       AgentStatus = "busy"       // 忙碌
	AgentStatusBreak      AgentStatus = "break"      // 小休
	AgentStatusAfterCall  AgentStatus = "after_call" // 话后处理
)

// Agent 座席
type Agent struct {
	ID           uuid.UUID              `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID              `gorm:"type:uuid;uniqueIndex" json:"user_id"`
	WorkID       string                 `gorm:"type:varchar(50);uniqueIndex" json:"work_id"`
	Name         string                 `gorm:"type:varchar(100)" json:"name"`
	Email        string                 `gorm:"type:varchar(100)" json:"email"`
	Phone        string                 `gorm:"type:varchar(20)" json:"phone"`
	Status       AgentStatus            `gorm:"type:varchar(20);index" json:"status"`
	Skills       []string               `gorm:"type:text[]" json:"skills"`
	Level        int                    `json:"level"`
	Extension    string                 `gorm:"type:varchar(20)" json:"extension"`
	LastLoginAt  *time.Time             `json:"last_login_at"`
	LastLogoutAt *time.Time             `json:"last_logout_at"`
	Metadata     map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

func (a *Agent) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// AgentSkillGroup 座席技能组
type AgentSkillGroup struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	AgentID     uuid.UUID `gorm:"type:uuid;index" json:"agent_id"`
	SkillGroup  string    `gorm:"type:varchar(50);index" json:"skill_group"`
	Priority    int       `json:"priority"`
	Proficiency int       `json:"proficiency"` // 熟练度 1-10
	CreatedAt   time.Time `json:"created_at"`
}

func (asg *AgentSkillGroup) BeforeCreate(tx *gorm.DB) error {
	if asg.ID == uuid.Nil {
		asg.ID = uuid.New()
	}
	return nil
}

// AgentStats 座席统计
type AgentStats struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	AgentID            uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_agent_date" json:"agent_id"`
	Date               string    `gorm:"type:varchar(10);uniqueIndex:idx_agent_date" json:"date"` // YYYY-MM-DD
	LoginDuration      int       `json:"login_duration"`      // 登录时长（秒）
	TalkDuration       int       `json:"talk_duration"`       // 通话时长（秒）
	IdleDuration       int       `json:"idle_duration"`       // 空闲时长（秒）
	BreakDuration      int       `json:"break_duration"`      // 小休时长（秒）
	AfterCallDuration  int       `json:"after_call_duration"` // 话后时长（秒）
	CallsHandled       int       `json:"calls_handled"`       // 处理呼叫数
	CallsAnswered      int       `json:"calls_answered"`      // 接听呼叫数
	CallsMissed        int       `json:"calls_missed"`        // 未接呼叫数
	CallsTransferred   int       `json:"calls_transferred"`   // 转接呼叫数
	AvgTalkTime        int       `json:"avg_talk_time"`       // 平均通话时长（秒）
	AvgWaitTime        int       `json:"avg_wait_time"`       // 平均等待时长（秒）
	AvgAfterCallTime   int       `json:"avg_after_call_time"` // 平均话后时长（秒）
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (as *AgentStats) BeforeCreate(tx *gorm.DB) error {
	if as.ID == uuid.Nil {
		as.ID = uuid.New()
	}
	return nil
}

// AgentSchedule 座席排班
type AgentSchedule struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	AgentID   uuid.UUID  `gorm:"type:uuid;index" json:"agent_id"`
	Date      string     `gorm:"type:varchar(10);index" json:"date"` // YYYY-MM-DD
	ShiftType string     `gorm:"type:varchar(20)" json:"shift_type"` // morning, afternoon, evening, night
	StartTime string     `gorm:"type:varchar(5)" json:"start_time"`  // HH:MM
	EndTime   string     `gorm:"type:varchar(5)" json:"end_time"`    // HH:MM
	BreakTime int        `json:"break_time"` // 休息时间（分钟）
	Status    string     `gorm:"type:varchar(20)" json:"status"` // scheduled, confirmed, cancelled
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (as *AgentSchedule) BeforeCreate(tx *gorm.DB) error {
	if as.ID == uuid.Nil {
		as.ID = uuid.New()
	}
	return nil
}
