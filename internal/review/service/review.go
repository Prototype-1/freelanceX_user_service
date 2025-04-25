package service

import (
	"time"
	"errors"
	"github.com/Prototype-1/freelanceX_user_service/internal/review/model"
	"github.com/Prototype-1/freelanceX_user_service/internal/review/repository"
	authRepo "github.com/Prototype-1/freelanceX_user_service/internal/auth/repository"
)

type ReviewService interface {
	SubmitReview(review *model.FreelancerReview) (*model.FreelancerReview, error)
	GetReviewsByFreelancerID(freelancerID string) ([]model.FreelancerReview, error)
}

type reviewService struct {
	repo repository.ReviewRepository
	userRepo authRepo.UserRepository
}

func NewReviewService(repo repository.ReviewRepository, userRepo authRepo.UserRepository) ReviewService {
	return &reviewService{
		repo : repo,
		userRepo: userRepo,
	}
}

var ErrUnauthorizedReview = errors.New("only clients are allowed to submit reviews")

func (s *reviewService) SubmitReview(review *model.FreelancerReview) (*model.FreelancerReview, error) {
	user, err := s.userRepo.GetUserByID(nil, review.ClientID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Role != "client" {
		return nil, ErrUnauthorizedReview 
	}

	review.CreatedAt = time.Now()
	if err := s.repo.Create(review); err != nil {
		return nil, err
	}
	return review, nil
}


func (s *reviewService) GetReviewsByFreelancerID(freelancerID string) ([]model.FreelancerReview, error) {
	return s.repo.GetByFreelancerID(freelancerID)
}
