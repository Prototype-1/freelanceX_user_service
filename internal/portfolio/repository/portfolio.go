package repository

import (
	"context"
	"github.com/Prototype-1/freelanceX_user_service/internal/portfolio/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, p *model.Portfolio) error
	GetByFreelancerID(ctx context.Context, freelancerID string) ([]*model.Portfolio, error)
	Delete(ctx context.Context, portfolioID string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, p *model.Portfolio) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *repository) GetByFreelancerID(ctx context.Context, freelancerID string) ([]*model.Portfolio, error) {
	var portfolios []*model.Portfolio
	err := r.db.WithContext(ctx).
		Where("freelancer_id = ?", freelancerID).
		Order("created_at desc").
		Find(&portfolios).Error
	return portfolios, err
}

func (r *repository) Delete(ctx context.Context, portfolioID string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", portfolioID).
		Delete(&model.Portfolio{}).Error
}
