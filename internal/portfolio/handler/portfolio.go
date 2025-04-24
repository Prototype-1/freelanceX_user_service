package handler

import (
	"context"
	"time"

	"github.com/Prototype-1/freelanceX_user_service/internal/portfolio/model"
	"github.com/Prototype-1/freelanceX_user_service/internal/portfolio/service"
	pb "github.com/Prototype-1/freelanceX_user_service/proto/portfolio"
)

type Handler struct {
	pb.UnimplementedPortfolioServiceServer
	service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreatePortfolio(ctx context.Context, req *pb.CreatePortfolioRequest) (*pb.CreatePortfolioResponse, error) {
	portfolio := &model.Portfolio{
		FreelancerID: req.FreelancerId,
		Title:        req.Title,
		Description:  req.Description,
		ImageURL:     req.ImageUrl,
		Link:         req.Link,
		CreatedAt:    time.Now(),
	}
	err := h.service.Create(ctx, portfolio)
	if err != nil {
		return nil, err
	}
	return &pb.CreatePortfolioResponse{Message: "Portfolio created"}, nil
}

func (h *Handler) GetPortfolio(ctx context.Context, req *pb.GetPortfolioRequest) (*pb.GetPortfolioResponse, error) {
	data, err := h.service.GetByFreelancerID(ctx, req.FreelancerId)
	if err != nil {
		return nil, err
	}
	var response []*pb.PortfolioItem
	for _, p := range data {
		response = append(response, &pb.PortfolioItem{
			Id:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			ImageUrl:    p.ImageURL,
			Link:        p.Link,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		})
	}
	return &pb.GetPortfolioResponse{Portfolio: response}, nil
}

func (h *Handler) DeletePortfolio(ctx context.Context, req *pb.DeletePortfolioRequest) (*pb.DeletePortfolioResponse, error) {
	err := h.service.Delete(ctx, req.PortfolioId)
	if err != nil {
		return nil, err
	}
	return &pb.DeletePortfolioResponse{Message: "Portfolio deleted"}, nil
}
