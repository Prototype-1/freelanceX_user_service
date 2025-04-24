package service

import (
	"context"

	"github.com/Prototype-1/freelanceX_user_service/internal/portfolio/model"
	"github.com/Prototype-1/freelanceX_user_service/internal/portfolio/repository"
)

type Service interface {
	Create(ctx context.Context, p *model.Portfolio) error
	GetByFreelancerID(ctx context.Context, freelancerID string) ([]*model.Portfolio, error)
	Delete(ctx context.Context, portfolioID string) error
}

type service struct {
	repo repository.Repository
}

func NewService(r repository.Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(ctx context.Context, p *model.Portfolio) error {
	return s.repo.Create(ctx, p)
}

func (s *service) GetByFreelancerID(ctx context.Context, freelancerID string) ([]*model.Portfolio, error) {
	return s.repo.GetByFreelancerID(ctx, freelancerID)
}

func (s *service) Delete(ctx context.Context, portfolioID string) error {
	return s.repo.Delete(ctx, portfolioID)
}
