package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// callRepository 呼叫仓储实现
type callRepository struct {
	db *gorm.DB
}

// NewCallRepository 创建呼叫仓储
func NewCallRepository(db *gorm.DB) CallRepository {
	return &callRepository{db: db}
}

// Create 创建呼叫记录
func (r *callRepository) Create(ctx context.Context, call *model.Call) error {
	return r.db.WithContext(ctx).Create(call).Error
}

// GetByID 根据ID获取呼叫
func (r *callRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Call, error) {
	var call model.Call
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&call).Error
	if err != nil {
		return nil, err
	}
	return &call, nil
}

// Update 更新呼叫信息
func (r *callRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Call{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除呼叫记录
func (r *callRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Call{}).Error
}

// List 查询呼叫列表
func (r *callRepository) List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*model.Call, error) {
	var calls []*model.Call
	
	query := r.db.WithContext(ctx).Model(&model.Call{})
	
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
	
	err := query.Find(&calls).Error
	return calls, err
}

// Count 统计呼叫数量
func (r *callRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	
	query := r.db.WithContext(ctx).Model(&model.Call{})
	
	// 应用过滤条件
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}
	
	err := query.Count(&count).Error
	return count, err
}
