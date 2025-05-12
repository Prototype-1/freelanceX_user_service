package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/Prototype-1/freelanceX_user_service/internal/profile/model"
	"github.com/Prototype-1/freelanceX_user_service/internal/profile/repository"
	pb "github.com/Prototype-1/freelanceX_user_service/proto/profile"
)

type Service interface {
	CreateOrUpdateProfile(ctx context.Context, req *pb.CreateProfileRequest) error
	GetProfile(ctx context.Context, userID string) (*model.FreelancerProfile, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo}
}

func (s *service) CreateOrUpdateProfile(ctx context.Context, req *pb.CreateProfileRequest) error {
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return err
	}

	profile := &model.FreelancerProfile{
		UserID:            userUUID,
		Title:             req.Title,
		Bio:               req.Bio,
		HourlyRate:        float64(req.HourlyRate),
		YearsOfExperience: int(req.YearsOfExperience),
		Skills:            pq.StringArray(req.Skills),
	Languages:         pq.StringArray(req.Languages),
	Certifications:    pq.StringArray(req.Certifications),
		Location:          req.Location,
		ResponseTime:      req.ResponseTime,
	}

	return s.repo.CreateOrUpdate(ctx, profile)
}

func (s *service) GetProfile(ctx context.Context, userID string) (*model.FreelancerProfile, error) {
	return s.repo.GetByUserID(ctx, userID)
}
