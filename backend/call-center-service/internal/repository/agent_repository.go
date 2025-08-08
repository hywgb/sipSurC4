package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// agentRepository 座席仓储实现
type agentRepository struct {
	db *gorm.DB
}

// NewAgentRepository 创建座席仓储
func NewAgentRepository(db *gorm.DB) AgentRepository {
	return &agentRepository{db: db}
}

// Create 创建座席
func (r *agentRepository) Create(ctx context.Context, agent *model.Agent) error {
	return r.db.WithContext(ctx).Create(agent).Error
}

// GetByID 根据ID获取座席
func (r *agentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Agent, error) {
	var agent model.Agent
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetByUserID 根据用户ID获取座席
func (r *agentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Agent, error) {
	var agent model.Agent
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetByWorkID 根据工号获取座席
func (r *agentRepository) GetByWorkID(ctx context.Context, workID string) (*model.Agent, error) {
	var agent model.Agent
	err := r.db.WithContext(ctx).Where("work_id = ?", workID).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// Update 更新座席信息
func (r *agentRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Agent{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除座席
func (r *agentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Agent{}).Error
}

// List 查询座席列表
func (r *agentRepository) List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*model.Agent, error) {
	var agents []*model.Agent
	
	query := r.db.WithContext(ctx).Model(&model.Agent{})
	
	// 应用过滤条件
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}
	
	// 应用分页
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	// 按创建时间倒序
	query = query.Order("created_at DESC")
	
	err := query.Find(&agents).Error
	return agents, err
}

// Count 统计座席数量
func (r *agentRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	
	query := r.db.WithContext(ctx).Model(&model.Agent{})
	
	// 应用过滤条件
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}
	
	err := query.Count(&count).Error
	return count, err
}

// UpdateStatus 更新座席状态
func (r *agentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.AgentStatus) error {
	return r.db.WithContext(ctx).Model(&model.Agent{}).
		Where("id = ?", id).
		Update("status", status).Error
}
