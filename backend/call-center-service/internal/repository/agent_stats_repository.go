package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// agentStatsRepository 座席统计仓储实现

type agentStatsRepository struct {
	db *gorm.DB
}

func NewAgentStatsRepository(db *gorm.DB) AgentStatsRepository {
	return &agentStatsRepository{db: db}
}

func (r *agentStatsRepository) Create(ctx context.Context, stats *model.AgentStats) error {
	return r.db.WithContext(ctx).Create(stats).Error
}

func (r *agentStatsRepository) GetByAgentIDAndDate(ctx context.Context, agentID uuid.UUID, date string) (*model.AgentStats, error) {
	var s model.AgentStats
	if err := r.db.WithContext(ctx).Where("agent_id = ? AND date = ?", agentID, date).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *agentStatsRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.AgentStats{}).Where("id = ?", id).Updates(updates).Error
}

func (r *agentStatsRepository) IncrementStats(ctx context.Context, agentID uuid.UUID, date string, field string, value int) error {
	return r.db.WithContext(ctx).Model(&model.AgentStats{}).
		Where("agent_id = ? AND date = ?", agentID, date).
		UpdateColumn(field, gorm.Expr(field+" + ?", value)).Error
}