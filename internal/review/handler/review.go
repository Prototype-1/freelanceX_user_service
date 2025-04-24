package handler

import (
	"context"
	"time"
	"github.com/Prototype-1/freelanceX_user_service/internal/review/model"
	"github.com/Prototype-1/freelanceX_user_service/internal/review/service"
	reviewPb "github.com/Prototype-1/freelanceX_user_service/proto/review"
)

type ReviewHandler struct {
	reviewPb.UnimplementedReviewServiceServer
	service service.ReviewService
}

func NewReviewHandler(s service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		service: s,
	}
}

func (h *ReviewHandler) SubmitReview(ctx context.Context, req *reviewPb.ReviewRequest) (*reviewPb.ReviewResponse, error) {
	review := &model.FreelancerReview{
		ProjectID:     req.GetProjectId(),
		FreelancerID:  req.GetFreelancerId(),
		ClientID:      req.GetClientId(),
		Rating:        int(req.GetRating()),
		Feedback:      req.GetFeedback(),
	}

	created, err := h.service.SubmitReview(review)
	if err != nil {
		return nil, err
	}

	return &reviewPb.ReviewResponse{
		Id:            created.ID,
		ProjectId:     created.ProjectID,
		FreelancerId:  created.FreelancerID,
		ClientId:      created.ClientID,
		Rating:        int32(created.Rating),
		Feedback:      created.Feedback,
		CreatedAt:     created.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *ReviewHandler) GetFreelancerReviews(ctx context.Context, req *reviewPb.GetReviewRequest) (*reviewPb.ReviewListResponse, error) {
	reviews, err := h.service.GetReviewsByFreelancerID(req.GetFreelancerId())
	if err != nil {
		return nil, err
	}

	var protoReviews []*reviewPb.ReviewResponse
	for _, r := range reviews {
		protoReviews = append(protoReviews, &reviewPb.ReviewResponse{
			Id:           r.ID,
			ProjectId:    r.ProjectID,
			FreelancerId: r.FreelancerID,
			ClientId:     r.ClientID,
			Rating:       int32(r.Rating),
			Feedback:     r.Feedback,
			CreatedAt:    r.CreatedAt.Format(time.RFC3339),
		})
	}

	return &reviewPb.ReviewListResponse{
		Reviews: protoReviews,
	}, nil
}
