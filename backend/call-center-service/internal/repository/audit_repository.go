package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// auditLogRepository 审计日志仓储实现

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository 创建审计日志仓储
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, logEntry *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(logEntry).Error
}

func (r *auditLogRepository) List(
	ctx context.Context,
	filter map[string]interface{},
	like map[string]string,
	rangeFilter map[string][2]interface{},
	offset, limit int,
) ([]*model.AuditLog, error) {
	var logs []*model.AuditLog
	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}
	for k, v := range like {
		query = query.Where(k+" ILIKE ?", "%"+v+"%")
	}
	for k, v := range rangeFilter {
		start := v[0]
		end := v[1]
		if start != nil {
			query = query.Where(k+" >= ?", start)
		}
		if end != nil {
			query = query.Where(k+" <= ?", end)
		}
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	query = query.Order("created_at DESC")
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *auditLogRepository) Count(
	ctx context.Context,
	filter map[string]interface{},
	like map[string]string,
	rangeFilter map[string][2]interface{},
) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}
	for k, v := range like {
		query = query.Where(k+" ILIKE ?", "%"+v+"%")
	}
	for k, v := range rangeFilter {
		start := v[0]
		end := v[1]
		if start != nil {
			query = query.Where(k+" >= ?", start)
		}
		if end != nil {
			query = query.Where(k+" <= ?", end)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}