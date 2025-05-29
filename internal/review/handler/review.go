package handler

import (
	"context"
	"time"
	"github.com/Prototype-1/freelanceX_user_service/internal/review/model"
	"github.com/Prototype-1/freelanceX_user_service/internal/review/service"
	reviewPb "github.com/Prototype-1/freelanceX_user_service/proto/review"
	"errors"
	"github.com/google/uuid"
	"github.com/Prototype-1/freelanceX_user_service/pkg/roles"
)

type ReviewHandler struct {
	reviewPb.UnimplementedReviewServiceServer
	service service.ReviewService
	roles role.Checker
}

func NewReviewHandler(s service.ReviewService, r role.Checker) *ReviewHandler {
	return &ReviewHandler{
		service: s,
		roles: r,
	}
}

var ErrPermissionDenied = errors.New("permission denied")

func (h *ReviewHandler) SubmitReview(ctx context.Context, req *reviewPb.ReviewRequest) (*reviewPb.ReviewResponse, error) {
	if !h.roles.HasRole(ctx, "client") {
		return nil, ErrPermissionDenied
	}

projectID, err := uuid.Parse(req.GetProjectId())
	if err != nil {
		return nil, err
	}
	freelancerID, err := uuid.Parse(req.GetFreelancerId())
	if err != nil {
		return nil, err
	}
	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, err
	}

	review := &model.FreelancerReview{
		ProjectID:    projectID,
		FreelancerID: freelancerID,
		ClientID:     clientID,
		Rating:       int(req.GetRating()),
		Feedback:     req.GetFeedback(),
	}

	created, err := h.service.SubmitReview(review)
	if err != nil {
		return nil, err
	}
return &reviewPb.ReviewResponse{
		Id:           created.ID.String(),
		ProjectId:    created.ProjectID.String(),
		FreelancerId: created.FreelancerID.String(),
		ClientId:     created.ClientID.String(),
		Rating:       int32(created.Rating),
		Feedback:     created.Feedback,
		CreatedAt:    created.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *ReviewHandler) GetFreelancerReviews(ctx context.Context, req *reviewPb.GetReviewRequest) (*reviewPb.ReviewListResponse, error) {
	if !h.roles.HasRole(ctx, "client", "admin") {
		return nil, ErrPermissionDenied
	}
	
	reviews, err := h.service.GetReviewsByFreelancerID(req.GetFreelancerId())
	if err != nil {
		return nil, err
	}

	var protoReviews []*reviewPb.ReviewResponse
for _, r := range reviews {
	protoReviews = append(protoReviews, &reviewPb.ReviewResponse{
		Id:           r.ID.String(),
		ProjectId:    r.ProjectID.String(),
		FreelancerId: r.FreelancerID.String(),
		ClientId:     r.ClientID.String(),
		Rating:       int32(r.Rating),
		Feedback:     r.Feedback,
		CreatedAt:    r.CreatedAt.Format(time.RFC3339),
	})
}
	return &reviewPb.ReviewListResponse{
		Reviews: protoReviews,
	}, nil
}
