package service

import (
	"context"
	"time"

	"github.com/cloudcallcenter/call-center-service/internal/model"
	"github.com/cloudcallcenter/call-center-service/internal/repository"
)

// auditService 审计日志服务实现

type auditService struct {
	repo repository.AuditLogRepository
}

func NewAuditService(repo repository.AuditLogRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) Log(ctx context.Context, entry *AuditLogInput) error {
	log := &model.AuditLog{
		TenantID:     entry.TenantID,
		ProjectID:    entry.ProjectID,
		ActorID:      entry.ActorID,
		ActorType:    entry.ActorType,
		Action:       entry.Action,
		ResourceType: entry.ResourceType,
		ResourceID:   entry.ResourceID,
		Level:        entry.Level,
		Method:       entry.Method,
		Path:         entry.Path,
		StatusCode:   entry.StatusCode,
		IP:           entry.IP,
		UserAgent:    entry.UserAgent,
		LatencyMs:    entry.LatencyMs,
		TraceID:      entry.TraceID,
		Metadata:     entry.Metadata,
		Error:        entry.Error,
		CreatedAt:    time.Now(),
	}
	return s.repo.Create(ctx, log)
}

func (s *auditService) List(ctx context.Context, filter *AuditLogFilter, page, pageSize int) ([]*model.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	f := map[string]interface{}{}
	like := map[string]string{}
	rangeFilter := map[string][2]interface{}{}

	if filter != nil {
		if filter.TenantID != nil {
			f["tenant_id"] = *filter.TenantID
		}
		if filter.ProjectID != nil {
			f["project_id"] = *filter.ProjectID
		}
		if filter.ActorID != nil {
			f["actor_id"] = *filter.ActorID
		}
		if filter.ActorType != nil {
			f["actor_type"] = *filter.ActorType
		}
		if filter.Action != nil {
			f["action"] = *filter.Action
		}
		if filter.ResourceType != nil {
			f["resource_type"] = *filter.ResourceType
		}
		if filter.ResourceID != nil {
			f["resource_id"] = *filter.ResourceID
		}
		if filter.Level != nil {
			f["level"] = *filter.Level
		}
		if filter.Method != nil {
			f["method"] = *filter.Method
		}
		if filter.PathLike != nil {
			like["path"] = *filter.PathLike
		}
		if filter.StatusCode != nil {
			f["status_code"] = *filter.StatusCode
		}
		if filter.StartTime != nil || filter.EndTime != nil {
			rangeFilter["created_at"] = [2]interface{}{nil, nil}
			if filter.StartTime != nil {
				rangeFilter["created_at"] = [2]interface{}{*filter.StartTime, rangeFilter["created_at"][1]}
			}
			if filter.EndTime != nil {
				rangeFilter["created_at"] = [2]interface{}{rangeFilter["created_at"][0], *filter.EndTime}
			}
		}
	}

	items, err := s.repo.List(ctx, f, like, rangeFilter, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(ctx, f, like, rangeFilter)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}