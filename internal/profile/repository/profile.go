package repository

import (
	"fmt"
	"context"
	"gorm.io/gorm"
	"github.com/Prototype-1/freelanceX_user_service/internal/profile/model"
	"github.com/google/uuid"
)

type Repository interface {
	CreateOrUpdate(ctx context.Context, profile *model.FreelancerProfile) error
	GetByUserID(ctx context.Context, userID string) (*model.FreelancerProfile, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateOrUpdate(ctx context.Context, profile *model.FreelancerProfile) error {
	var existing model.FreelancerProfile
	tx := r.db.WithContext(ctx).Where("user_id = ?", profile.UserID).First(&existing)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return tx.Error
	}

	if existing.ID != uuid.Nil {
		profile.ID = existing.ID
		return r.db.WithContext(ctx).Save(profile).Error
	}

	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *repository) GetByUserID(ctx context.Context, userID string) (*model.FreelancerProfile, error) {
_, err := uuid.Parse(userID)
if err != nil {
	return nil, fmt.Errorf("invalid UUID format for user_id: %v", err)
}
	var profile model.FreelancerProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}
