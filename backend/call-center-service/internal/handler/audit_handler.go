package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/cloudcallcenter/call-center-service/internal/service"
)

// AuditHandler 审计日志处理器

type AuditHandler struct {
	service service.AuditService
}

func NewAuditHandler(s service.AuditService) *AuditHandler {
	return &AuditHandler{service: s}
}

// ListAudits 查询审计日志
func (h *AuditHandler) ListAudits(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := &service.AuditLogFilter{}

	if tenantID := c.Query("tenant_id"); tenantID != "" {
		if id, err := uuid.Parse(tenantID); err == nil {
			filter.TenantID = &id
		}
	}
	if projectID := c.Query("project_id"); projectID != "" {
		if id, err := uuid.Parse(projectID); err == nil {
			filter.ProjectID = &id
		}
	}
	if actorID := c.Query("actor_id"); actorID != "" {
		if id, err := uuid.Parse(actorID); err == nil {
			filter.ActorID = &id
		}
	}
	if v := c.Query("actor_type"); v != "" {
		filter.ActorType = &v
	}
	if v := c.Query("action"); v != "" {
		filter.Action = &v
	}
	if v := c.Query("resource_type"); v != "" {
		filter.ResourceType = &v
	}
	if v := c.Query("resource_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.ResourceID = &id
		}
	}
	if v := c.Query("level"); v != "" {
		filter.Level = &v
	}
	if v := c.Query("method"); v != "" {
		filter.Method = &v
	}
	if v := c.Query("path_like"); v != "" {
		filter.PathLike = &v
	}
	if v := c.Query("status_code"); v != "" {
		if code, err := strconv.Atoi(v); err == nil {
			filter.StatusCode = &code
		}
	}
	if v := c.Query("start_time"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.StartTime = &t
		}
	}
	if v := c.Query("end_time"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.EndTime = &t
		}
	}

	items, total, err := h.service.List(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}