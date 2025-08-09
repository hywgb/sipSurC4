package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// agentScheduleRepository 座席排班仓储实现

type agentScheduleRepository struct {
	db *gorm.DB
}

func NewAgentScheduleRepository(db *gorm.DB) AgentScheduleRepository {
	return &agentScheduleRepository{db: db}
}

func (r *agentScheduleRepository) Create(ctx context.Context, schedule *model.AgentSchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *agentScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.AgentSchedule, error) {
	var s model.AgentSchedule
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *agentScheduleRepository) GetByAgentIDAndDateRange(ctx context.Context, agentID uuid.UUID, startDate, endDate string) ([]*model.AgentSchedule, error) {
	var schedules []*model.AgentSchedule
	if err := r.db.WithContext(ctx).
		Where("agent_id = ? AND date >= ? AND date <= ?", agentID, startDate, endDate).
		Order("date ASC, start_time ASC").
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *agentScheduleRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.AgentSchedule{}).Where("id = ?", id).Updates(updates).Error
}

func (r *agentScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.AgentSchedule{}).Error
}