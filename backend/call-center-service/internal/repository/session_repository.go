package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// sessionRepository 会话仓储实现
type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository 创建会话仓储
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

// Create 创建会话
func (r *sessionRepository) Create(ctx context.Context, session *model.CallSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID 根据ID获取会话
func (r *sessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CallSession, error) {
	var session model.CallSession
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Update 更新会话信息
func (r *sessionRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.CallSession{}).Where("id = ?", id).Updates(updates).Error
}

// List 查询会话列表
func (r *sessionRepository) List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*model.CallSession, error) {
	var sessions []*model.CallSession
	
	query := r.db.WithContext(ctx).Model(&model.CallSession{})
	
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
	
	err := query.Find(&sessions).Error
	return sessions, err
}

// GetByCallID 根据呼叫ID获取会话
func (r *sessionRepository) GetByCallID(ctx context.Context, callID uuid.UUID) ([]*model.CallSession, error) {
	var sessions []*model.CallSession
	err := r.db.WithContext(ctx).
		Where("call_id = ?", callID).
		Order("join_time DESC").
		Find(&sessions).Error
	return sessions, err
}

// GetActiveSessions 获取活跃会话
func (r *sessionRepository) GetActiveSessions(ctx context.Context) ([]*model.CallSession, error) {
	var sessions []*model.CallSession
	err := r.db.WithContext(ctx).
		Where("leave_time IS NULL").
		Order("join_time DESC").
		Find(&sessions).Error
	return sessions, err
}

// Count 统计会话数量
func (r *sessionRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.CallSession{})
	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
