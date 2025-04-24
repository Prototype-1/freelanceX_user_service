package repository

import (
	"github.com/Prototype-1/freelanceX_user_service/internal/review/model"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *model.FreelancerReview) error
	GetByFreelancerID(freelancerID string) ([]model.FreelancerReview, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db}
}

func (r *reviewRepository) Create(review *model.FreelancerReview) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetByFreelancerID(freelancerID string) ([]model.FreelancerReview, error) {
	var reviews []model.FreelancerReview
	err := r.db.Where("freelancer_id = ?", freelancerID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}
