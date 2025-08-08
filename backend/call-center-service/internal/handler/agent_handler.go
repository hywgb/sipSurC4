package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/cloudcallcenter/call-center-service/internal/model"
	"github.com/cloudcallcenter/call-center-service/internal/service"
)

// AgentHandler 座席处理器
type AgentHandler struct {
	agentService service.AgentService
}

// NewAgentHandler 创建座席处理器
func NewAgentHandler(agentService service.AgentService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
	}
}

// ListAgents 查询座席列表
func (h *AgentHandler) ListAgents(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 构建过滤条件
	filter := &service.AgentFilter{}
	
	if status := c.Query("status"); status != "" {
		s := model.AgentStatus(status)
		filter.Status = &s
	}
	
	if level := c.Query("level"); level != "" {
		if l, err := strconv.Atoi(level); err == nil {
			filter.Level = &l
		}
	}

	// 查询数据
	agents, total, err := h.agentService.ListAgents(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      agents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetAgent 获取座席信息
func (h *AgentHandler) GetAgent(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent ID"})
		return
	}

	agent, err := h.agentService.GetAgent(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// UpdateAgentStatus 更新座席状态
func (h *AgentHandler) UpdateAgentStatus(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := model.AgentStatus(req.Status)
	if err := h.agentService.UpdateAgentStatus(c.Request.Context(), agentID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "agent status updated successfully"})
}

// GetAgentStats 获取座席统计
func (h *AgentHandler) GetAgentStats(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent ID"})
		return
	}

	date := c.DefaultQuery("date", "")
	if date == "" {
		// 默认获取今天的统计
		date = time.Now().Format("2006-01-02")
	}

	stats, err := h.agentService.GetAgentStats(c.Request.Context(), agentID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Login 座席登录
func (h *AgentHandler) Login(c *gin.Context) {
	var req struct {
		AgentID string `json:"agent_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentID, err := uuid.Parse(req.AgentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent ID"})
		return
	}

	if err := h.agentService.Login(c.Request.Context(), agentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

// Logout 座席登出
func (h *AgentHandler) Logout(c *gin.Context) {
	var req struct {
		AgentID string `json:"agent_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentID, err := uuid.Parse(req.AgentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent ID"})
		return
	}

	if err := h.agentService.Logout(c.Request.Context(), agentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
