package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/cloudcallcenter/call-center-service/internal/service"
)

// CallHandler 呼叫处理器
type CallHandler struct {
	callService service.CallService
}

// NewCallHandler 创建呼叫处理器
func NewCallHandler(callService service.CallService) *CallHandler {
	return &CallHandler{
		callService: callService,
	}
}

// CreateCall 创建呼叫
func (h *CallHandler) CreateCall(c *gin.Context) {
	var req service.CreateCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	call, err := h.callService.CreateCall(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, call)
}

// GetCall 获取呼叫信息
func (h *CallHandler) GetCall(c *gin.Context) {
	callID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call ID"})
		return
	}

	call, err := h.callService.GetCall(c.Request.Context(), callID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "call not found"})
		return
	}

	c.JSON(http.StatusOK, call)
}

// UpdateCall 更新呼叫信息
func (h *CallHandler) UpdateCall(c *gin.Context) {
	callID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.callService.UpdateCall(c.Request.Context(), callID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "call updated successfully"})
}

// HangupCall 挂断呼叫
func (h *CallHandler) HangupCall(c *gin.Context) {
	callID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call ID"})
		return
	}

	if err := h.callService.HangupCall(c.Request.Context(), callID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "call hung up successfully"})
}

// GetRecording 获取录音
func (h *CallHandler) GetRecording(c *gin.Context) {
	callID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call ID"})
		return
	}

	recording, err := h.callService.GetCallRecording(c.Request.Context(), callID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recording not found"})
		return
	}

	c.JSON(http.StatusOK, recording)
}

// ListCalls 查询呼叫列表
func (h *CallHandler) ListCalls(c *gin.Context) {
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
	filter := &service.CallFilter{}
	
	if projectID := c.Query("project_id"); projectID != "" {
		if id, err := uuid.Parse(projectID); err == nil {
			filter.ProjectID = &id
		}
	}
	
	if status := c.Query("status"); status != "" {
		s := service.CallStatus(status)
		filter.Status = (*service.CallStatus)(&s)
	}

	// 查询数据
	calls, total, err := h.callService.ListCalls(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      calls,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetCallStats 获取呼叫统计
func (h *CallHandler) GetCallStats(c *gin.Context) {
	var tr service.TimeRange
	if v := c.Query("start"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			tr.Start = t
		}
	}
	if v := c.Query("end"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			tr.End = t
		}
	}
	stats, err := h.callService.GetCallStats(c.Request.Context(), tr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
