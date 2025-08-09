package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cloudcallcenter/call-center-service/internal/model"
)

// recordingRepository 录音仓储实现

type recordingRepository struct {
	db *gorm.DB
}

// NewRecordingRepository 创建录音仓储
func NewRecordingRepository(db *gorm.DB) RecordingRepository {
	return &recordingRepository{db: db}
}

func (r *recordingRepository) Create(ctx context.Context, recording *model.Recording) error {
	return r.db.WithContext(ctx).Create(recording).Error
}

func (r *recordingRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Recording, error) {
	var rec model.Recording
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *recordingRepository) GetByCallID(ctx context.Context, callID uuid.UUID) (*model.Recording, error) {
	var rec model.Recording
	if err := r.db.WithContext(ctx).Where("call_id = ?", callID).First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *recordingRepository) List(ctx context.Context, filter map[string]interface{}, offset, limit int) ([]*model.Recording, error) {
	var recs []*model.Recording
	q := r.db.WithContext(ctx).Model(&model.Recording{})
	for k, v := range filter {
		q = q.Where(k+" = ?", v)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	q = q.Order("created_at DESC")
	if err := q.Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}