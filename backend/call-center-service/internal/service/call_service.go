package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	
	"github.com/cloudcallcenter/call-center-service/internal/model"
	"github.com/cloudcallcenter/call-center-service/internal/repository"
	"github.com/cloudcallcenter/call-center-service/pkg/dialer"
	"github.com/cloudcallcenter/call-center-service/pkg/router"
)

// callService 呼叫服务实现
type callService struct {
	callRepo repository.CallRepository
	dialer   dialer.Dialer
	router   router.Router
	logger   *logrus.Logger
}

// NewCallService 创建呼叫服务
func NewCallService(
	callRepo repository.CallRepository,
	dialer dialer.Dialer,
	router router.Router,
) CallService {
	return &callService{
		callRepo: callRepo,
		dialer:   dialer,
		router:   router,
		logger:   logrus.New(),
	}
}

// CreateCall 创建呼叫
func (s *callService) CreateCall(ctx context.Context, req *CreateCallRequest) (*model.Call, error) {
	// 验证请求
	if err := s.validateCreateCallRequest(req); err != nil {
		return nil, err
	}

	// 创建呼叫记录
	call := &model.Call{
		ID:            uuid.New(),
		ProjectID:     req.ProjectID,
		CampaignID:    req.CampaignID,
		CustomerPhone: req.CustomerPhone,
		CallerID:      req.CallerID,
		Direction:     req.Direction,
		Status:        model.CallStatusPending,
		DialMode:      req.DialMode,
		Metadata:      req.Metadata,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 保存到数据库
	if err := s.callRepo.Create(ctx, call); err != nil {
		return nil, err
	}

	// 如果是外呼，需要路由和拨号
	if req.Direction == model.CallDirectionOutbound {
		// 路由到座席
		routeReq := &router.RouteRequest{
			CallID:         call.ID,
			RequiredSkills: req.Skills,
			Priority:       req.Priority,
			Metadata:       req.Metadata,
		}

		routeResult, err := s.router.Route(ctx, routeReq)
		if err != nil {
			s.logger.Errorf("Failed to route call %s: %v", call.ID, err)
			call.Status = model.CallStatusFailed
			s.callRepo.Update(ctx, call.ID, map[string]interface{}{
				"status": call.Status,
			})
			return nil, err
		}

		// 更新座席ID
		call.AgentID = &routeResult.AgentID
		s.callRepo.Update(ctx, call.ID, map[string]interface{}{
			"agent_id": call.AgentID,
		})

		// 执行拨号
		dialReq := &dialer.DialRequest{
			ID:            call.ID,
			ProjectID:     req.ProjectID,
			CampaignID:    req.CampaignID,
			CustomerPhone: req.CustomerPhone,
			CallerID:      req.CallerID,
			Mode:          dialer.DialMode(req.DialMode),
			Priority:      req.Priority,
			Metadata:      req.Metadata,
		}

		dialResult, err := s.dialer.Dial(ctx, dialReq)
		if err != nil {
			s.logger.Errorf("Failed to dial call %s: %v", call.ID, err)
			call.Status = model.CallStatusFailed
			s.callRepo.Update(ctx, call.ID, map[string]interface{}{
				"status": call.Status,
			})
			return nil, err
		}

		// 更新呼叫状态
		call.Status = model.CallStatusDialing
		call.StartTime = &[]time.Time{time.Now()}[0]
		s.callRepo.Update(ctx, call.ID, map[string]interface{}{
			"status":     call.Status,
			"start_time": call.StartTime,
			"call_sid":   dialResult.CallID.String(),
		})
	}

	s.logger.Infof("Created call %s for customer %s", call.ID, call.CustomerPhone)
	return call, nil
}

// GetCall 获取呼叫信息
func (s *callService) GetCall(ctx context.Context, callID uuid.UUID) (*model.Call, error) {
	return s.callRepo.GetByID(ctx, callID)
}

// UpdateCall 更新呼叫信息
func (s *callService) UpdateCall(ctx context.Context, callID uuid.UUID, updates map[string]interface{}) error {
	// 添加更新时间
	updates["updated_at"] = time.Now()
	return s.callRepo.Update(ctx, callID, updates)
}

// ListCalls 查询呼叫列表
func (s *callService) ListCalls(ctx context.Context, filter *CallFilter, page, pageSize int) ([]*model.Call, int64, error) {
	// 构建查询条件
	query := s.buildCallQuery(filter)
	
	// 计算偏移量
	offset := (page - 1) * pageSize
	
	// 查询数据
	calls, err := s.callRepo.List(ctx, query, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	// 查询总数
	total, err := s.callRepo.Count(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	
	return calls, total, nil
}

// HangupCall 挂断呼叫
func (s *callService) HangupCall(ctx context.Context, callID uuid.UUID) error {
	// 获取呼叫信息
	call, err := s.callRepo.GetByID(ctx, callID)
	if err != nil {
		return err
	}

	// 检查呼叫状态
	if call.Status != model.CallStatusConnected && call.Status != model.CallStatusRinging {
		return errors.New("call is not in a state that can be hung up")
	}

	// 取消拨号（如果还在拨号中）
	if call.Status == model.CallStatusDialing || call.Status == model.CallStatusRinging {
		if err := s.dialer.CancelDial(ctx, callID); err != nil {
			s.logger.Errorf("Failed to cancel dial for call %s: %v", callID, err)
		}
	}

	// 更新呼叫状态
	endTime := time.Now()
	duration := 0
	if call.ConnectTime != nil {
		duration = int(endTime.Sub(*call.ConnectTime).Seconds())
	}

	updates := map[string]interface{}{
		"status":   model.CallStatusHangup,
		"end_time": endTime,
		"duration": duration,
	}

	if err := s.callRepo.Update(ctx, callID, updates); err != nil {
		return err
	}

	s.logger.Infof("Hung up call %s", callID)
	return nil
}

// GetCallRecording 获取呼叫录音
func (s *callService) GetCallRecording(ctx context.Context, callID uuid.UUID) (*model.Recording, error) {
	// TODO: 实现录音查询
	return nil, errors.New("not implemented")
}

// GetCallStats 获取呼叫统计
func (s *callService) GetCallStats(ctx context.Context, timeRange TimeRange) (*CallStats, error) {
	// TODO: 实现统计查询
	stats := &CallStats{
		TotalCalls:      100,
		AnsweredCalls:   80,
		MissedCalls:     15,
		FailedCalls:     5,
		AvgCallDuration: 180.5,
		AvgWaitTime:     15.3,
		AnswerRate:      0.8,
		AbandonRate:     0.02,
	}
	return stats, nil
}

// validateCreateCallRequest 验证创建呼叫请求
func (s *callService) validateCreateCallRequest(req *CreateCallRequest) error {
	if req.CustomerPhone == "" {
		return errors.New("customer phone is required")
	}
	if req.ProjectID == uuid.Nil {
		return errors.New("project ID is required")
	}
	if req.CampaignID == uuid.Nil {
		return errors.New("campaign ID is required")
	}
	return nil
}

// buildCallQuery 构建查询条件
func (s *callService) buildCallQuery(filter *CallFilter) map[string]interface{} {
	query := make(map[string]interface{})
	
	if filter.ProjectID != nil {
		query["project_id"] = *filter.ProjectID
	}
	if filter.CampaignID != nil {
		query["campaign_id"] = *filter.CampaignID
	}
	if filter.AgentID != nil {
		query["agent_id"] = *filter.AgentID
	}
	if filter.Status != nil {
		query["status"] = *filter.Status
	}
	if filter.Direction != nil {
		query["direction"] = *filter.Direction
	}
	if filter.Phone != nil {
		query["customer_phone"] = *filter.Phone
	}
	
	return query
}
