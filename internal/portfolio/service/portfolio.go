package service

import (
	"context"
	"github.com/Prototype-1/freelanceX_user_service/internal/portfolio/repository"
	"github.com/Prototype-1/freelanceX_user_service/proto/portfolio"
)

type PortfolioService struct {
	Repo repository.PortfolioRepository
}

func NewPortfolioService(repo repository.PortfolioRepository) *PortfolioService {
	return &PortfolioService{Repo: repo}
}

func (s *PortfolioService) CreatePortfolio(ctx context.Context, req *portfolio.CreatePortfolioRequest) (*portfolio.CreatePortfolioResponse, error) {
	err := s.Repo.CreatePortfolio(ctx, req)
	if err != nil {
		return nil, err
	}
	return &portfolio.CreatePortfolioResponse{
		Message: "Portfolio created successfully",
	}, nil
}

func (s *PortfolioService) GetPortfolio(ctx context.Context, req *portfolio.GetPortfolioRequest) (*portfolio.GetPortfolioResponse, error) {
	portfolios, err := s.Repo.GetPortfolio(ctx, req.FreelancerId)
	if err != nil {
		return nil, err
	}
	return &portfolio.GetPortfolioResponse{Portfolio: portfolios}, nil
}

func (s *PortfolioService) DeletePortfolio(ctx context.Context, req *portfolio.DeletePortfolioRequest) (*portfolio.DeletePortfolioResponse, error) {
	err := s.Repo.DeletePortfolio(ctx, req.PortfolioId)
	if err != nil {
		return nil, err
	}
	return &portfolio.DeletePortfolioResponse{
		Message: "Portfolio item deleted successfully",
	}, nil
}
