package handler

import (
	"context"
	"time"
 "github.com/Prototype-1/freelanceX_user_service/pkg/roles"
	"github.com/Prototype-1/freelanceX_user_service/internal/portfolio/model"
	"github.com/Prototype-1/freelanceX_user_service/internal/portfolio/service"
	pb "github.com/Prototype-1/freelanceX_user_service/proto/portfolio"
	"errors"
)

type Handler struct {
	pb.UnimplementedPortfolioServiceServer
	service service.Service
	roles   role.Checker
}

func NewHandler(s service.Service, r role.Checker) *Handler {
	return &Handler{service: s, roles: r}
}

var ErrPermissionDenied = errors.New("permission denied")

func (h *Handler) CreatePortfolio(ctx context.Context, req *pb.CreatePortfolioRequest) (*pb.CreatePortfolioResponse, error) {
	if !h.roles.HasRole(ctx, "freelancer") {
		return nil, ErrPermissionDenied
	}
	
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
	if !h.roles.HasRole(ctx, "freelancer", "client", "admin") {
		return nil, ErrPermissionDenied
	}

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
	if !h.roles.HasRole(ctx, "freelancer", "admin") {
		return nil, ErrPermissionDenied
	}
	
	err := h.service.Delete(ctx, req.PortfolioId)
	if err != nil {
		return nil, err
	}
	return &pb.DeletePortfolioResponse{Message: "Portfolio deleted"}, nil
}
