package service

import (
	"time"

	"github.com/Prototype-1/freelanceX_user_service/internal/review/model"
	"github.com/Prototype-1/freelanceX_user_service/internal/review/repository"
)

type ReviewService interface {
	SubmitReview(review *model.FreelancerReview) (*model.FreelancerReview, error)
	GetReviewsByFreelancerID(freelancerID string) ([]model.FreelancerReview, error)
}

type reviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) ReviewService {
	return &reviewService{repo}
}

func (s *reviewService) SubmitReview(review *model.FreelancerReview) (*model.FreelancerReview, error) {
	review.CreatedAt = time.Now()
	if err := s.repo.Create(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *reviewService) GetReviewsByFreelancerID(freelancerID string) ([]model.FreelancerReview, error) {
	return s.repo.GetByFreelancerID(freelancerID)
}
