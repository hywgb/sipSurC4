package dialer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrDialerBusy         = errors.New("dialer is busy")
	ErrNoAvailableAgent   = errors.New("no available agent")
)

// DialRequest 拨号请求
type DialRequest struct {
	ID            uuid.UUID
	ProjectID     uuid.UUID
	CampaignID    uuid.UUID
	CustomerPhone string
	CallerID      string
	Mode          DialMode
	Priority      int
	Metadata      map[string]interface{}
}

// DialResult 拨号结果
type DialResult struct {
	CallID   uuid.UUID
	Status   DialStatus
	Duration time.Duration
	Error    error
}

// DialStatus 拨号状态
type DialStatus string

const (
	DialStatusSuccess   DialStatus = "success"
	DialStatusFailed    DialStatus = "failed"
	DialStatusNoAnswer  DialStatus = "no_answer"
	DialStatusBusy      DialStatus = "busy"
	DialStatusCancelled DialStatus = "cancelled"
)

// DialMode 拨号模式
type DialMode string

const (
	DialModePredictive DialMode = "predictive"
	DialModePreview    DialMode = "preview"
	DialModeAuto       DialMode = "auto"
	DialModeManual     DialMode = "manual"
)

// Dialer 拨号器接口
type Dialer interface {
	Dial(ctx context.Context, req *DialRequest) (*DialResult, error)
	CancelDial(ctx context.Context, callID uuid.UUID) error
	GetDialerStats() *DialerStats
	SetDialRatio(ratio float64)
	Start(ctx context.Context) error
	Stop() error
}

// DialerStats 拨号器统计
type DialerStats struct {
	TotalCalls      int64
	SuccessfulCalls int64
	FailedCalls     int64
	ActiveCalls     int64
	AnswerRate      float64
	AbandonRate     float64
	AvgWaitTime     time.Duration
}

// PredictiveDialer 预测式拨号器
type PredictiveDialer struct {
	mu           sync.RWMutex
	dialRatio    float64
	maxConcurrent int
	activeCalls   map[uuid.UUID]*DialRequest
	stats        *DialerStats
	logger       *logrus.Logger
}

// NewDialer 创建拨号器
func NewDialer() Dialer {
	return &PredictiveDialer{
		dialRatio:     1.5,
		maxConcurrent: 1000,
		activeCalls:   make(map[uuid.UUID]*DialRequest),
		stats:         &DialerStats{},
		logger:        logrus.New(),
	}
}

// Dial 执行拨号
func (d *PredictiveDialer) Dial(ctx context.Context, req *DialRequest) (*DialResult, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 检查并发限制
	if len(d.activeCalls) >= d.maxConcurrent {
		return nil, ErrDialerBusy
	}

	// 验证电话号码
	if err := validatePhoneNumber(req.CustomerPhone); err != nil {
		return nil, ErrInvalidPhoneNumber
	}

	// 记录活跃呼叫
	d.activeCalls[req.ID] = req
	d.stats.ActiveCalls = int64(len(d.activeCalls))

	// 异步执行拨号
	go d.executeDial(ctx, req)

	return &DialResult{
		CallID: req.ID,
		Status: DialStatusSuccess,
	}, nil
}

// executeDial 执行实际的拨号操作
func (d *PredictiveDialer) executeDial(ctx context.Context, req *DialRequest) {
	d.logger.Infof("Dialing %s with mode %s", req.CustomerPhone, req.Mode)

	// 模拟拨号延迟
	time.Sleep(2 * time.Second)

	// 更新统计
	d.mu.Lock()
	delete(d.activeCalls, req.ID)
	d.stats.TotalCalls++
	d.stats.SuccessfulCalls++
	d.stats.ActiveCalls = int64(len(d.activeCalls))
	d.mu.Unlock()

	// TODO: 实际的CTI集成
}

// CancelDial 取消拨号
func (d *PredictiveDialer) CancelDial(ctx context.Context, callID uuid.UUID) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.activeCalls[callID]; !exists {
		return fmt.Errorf("call %s not found", callID)
	}

	delete(d.activeCalls, callID)
	d.stats.ActiveCalls = int64(len(d.activeCalls))

	return nil
}

// GetDialerStats 获取拨号器统计
func (d *PredictiveDialer) GetDialerStats() *DialerStats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	stats := *d.stats
	if stats.TotalCalls > 0 {
		stats.AnswerRate = float64(stats.SuccessfulCalls) / float64(stats.TotalCalls)
	}

	return &stats
}

// SetDialRatio 设置拨号比
func (d *PredictiveDialer) SetDialRatio(ratio float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.dialRatio = ratio
}

// Start 启动拨号器
func (d *PredictiveDialer) Start(ctx context.Context) error {
	d.logger.Info("Starting predictive dialer")
	
	// 启动监控协程
	go d.monitor(ctx)
	
	return nil
}

// Stop 停止拨号器
func (d *PredictiveDialer) Stop() error {
	d.logger.Info("Stopping predictive dialer")
	
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// 清理所有活跃呼叫
	for callID := range d.activeCalls {
		delete(d.activeCalls, callID)
	}
	
	return nil
}

// monitor 监控拨号器状态
func (d *PredictiveDialer) monitor(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			stats := d.GetDialerStats()
			d.logger.WithFields(logrus.Fields{
				"total_calls":   stats.TotalCalls,
				"active_calls":  stats.ActiveCalls,
				"answer_rate":   stats.AnswerRate,
				"abandon_rate":  stats.AbandonRate,
			}).Info("Dialer stats")

			// 动态调整拨号比
			d.adjustDialRatio(stats)
		}
	}
}

// adjustDialRatio 动态调整拨号比
func (d *PredictiveDialer) adjustDialRatio(stats *DialerStats) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 如果放弃率超过3%，降低拨号比
	if stats.AbandonRate > 0.03 {
		d.dialRatio *= 0.9
		d.logger.Warnf("High abandon rate %.2f%%, reducing dial ratio to %.2f", 
			stats.AbandonRate*100, d.dialRatio)
	} else if stats.AbandonRate < 0.01 && d.dialRatio < 3.0 {
		// 如果放弃率很低，可以适当提高拨号比
		d.dialRatio *= 1.1
		d.logger.Infof("Low abandon rate %.2f%%, increasing dial ratio to %.2f", 
			stats.AbandonRate*100, d.dialRatio)
	}

	// 限制拨号比范围
	if d.dialRatio < 1.0 {
		d.dialRatio = 1.0
	} else if d.dialRatio > 3.0 {
		d.dialRatio = 3.0
	}
}

// validatePhoneNumber 验证电话号码
func validatePhoneNumber(phone string) error {
	// 简单的电话号码验证
	if len(phone) < 7 || len(phone) > 15 {
		return ErrInvalidPhoneNumber
	}
	return nil
}
