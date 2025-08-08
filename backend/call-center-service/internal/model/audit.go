package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog 审计日志
// 记录系统操作以满足合规与追踪
// ResourceType 示例：call, agent, session, recording, system
// Action 示例：create, update, delete, get, login, logout, hangup, route
// ActorType 示例：agent, admin, system, api
// Level 示例：info, warn, error
// Note: 采用 JSONB 存储 metadata 便于扩展
// Indexes: 依据常用查询维度建立
// Timezone: 使用服务器时区（推荐 UTC）

type AuditLog struct {
	ID           uuid.UUID              `gorm:"type:uuid;primary_key" json:"id"`
	TenantID     *uuid.UUID             `gorm:"type:uuid;index" json:"tenant_id"`
	ProjectID    *uuid.UUID             `gorm:"type:uuid;index" json:"project_id"`
	ActorID      *uuid.UUID             `gorm:"type:uuid;index" json:"actor_id"`
	ActorType    string                 `gorm:"type:varchar(20);index" json:"actor_type"`
	Action       string                 `gorm:"type:varchar(50);index" json:"action"`
	ResourceType string                 `gorm:"type:varchar(50);index" json:"resource_type"`
	ResourceID   *uuid.UUID             `gorm:"type:uuid;index" json:"resource_id"`
	Level        string                 `gorm:"type:varchar(10);index" json:"level"`
	Method       string                 `gorm:"type:varchar(10);index" json:"method"`
	Path         string                 `gorm:"type:varchar(255);index" json:"path"`
	StatusCode   int                    `json:"status_code"`
	IP           string                 `gorm:"type:varchar(64)" json:"ip"`
	UserAgent    string                 `gorm:"type:varchar(255)" json:"user_agent"`
	LatencyMs    int64                  `json:"latency_ms"`
	TraceID      string                 `gorm:"type:varchar(64);index" json:"trace_id"`
	Metadata     map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	Error        string                 `gorm:"type:text" json:"error"`
	CreatedAt    time.Time              `json:"created_at"`
}

func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	return nil
}