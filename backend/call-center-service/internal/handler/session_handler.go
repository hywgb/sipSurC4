package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/cloudcallcenter/call-center-service/internal/service"
)

// SessionHandler 会话处理器
type SessionHandler struct {
	sessionService service.SessionService
}

// NewSessionHandler 创建会话处理器
func NewSessionHandler(sessionService service.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

// ListSessions 查询会话列表
func (h *SessionHandler) ListSessions(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filter := &service.SessionFilter{}
	if v := c.Query("call_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.CallID = &id
		}
	}
	if v := c.Query("agent_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.AgentID = &id
		}
	}
	if v := c.Query("state"); v != "" {
		filter.State = &v
	}

	sessions, total, err := h.sessionService.ListSessions(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      sessions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetSession 获取会话信息
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetActiveSessions 获取活跃会话
func (h *SessionHandler) GetActiveSessions(c *gin.Context) {
	sessions, err := h.sessionService.GetActiveSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  sessions,
		"total": len(sessions),
	})
}
