package handler

import (
	"context"
	"github.com/Prototype-1/freelanceX_user_service/internal/profile/service"
	"github.com/Prototype-1/freelanceX_user_service/proto/profile"
	"github.com/Prototype-1/freelanceX_user_service/pkg/roles"
	"errors"
)

type Handler struct {
	profile.UnimplementedProfileServiceServer
	service service.Service
	roles   role.Checker
}

func NewHandler(s service.Service, r role.Checker) *Handler {
	return &Handler{service: s, roles: r}
}

var ErrPermissionDenied = errors.New("permission denied")

func (h *Handler) CreateProfile(ctx context.Context, req *profile.CreateProfileRequest) (*profile.CreateProfileResponse, error) {
	if !h.roles.HasRole(ctx, "freelancer") {
		return nil, ErrPermissionDenied
	}

	err := h.service.CreateOrUpdateProfile(ctx, req)
	if err != nil {
		return nil, err
	}
	return &profile.CreateProfileResponse{Message: "Profile created/updated successfully"}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, req *profile.UpdateProfileRequest) (*profile.UpdateProfileResponse, error) {
	if !h.roles.HasRole(ctx, "freelancer") {
		return nil, ErrPermissionDenied
	}

	err := h.service.CreateOrUpdateProfile(ctx, &profile.CreateProfileRequest{
		UserId:            req.UserId,
		Title:              req.Title,
		Bio:                req.Bio,
		HourlyRate:        req.HourlyRate,
		YearsOfExperience: req.YearsOfExperience,
		Skills:             req.Skills,
		Languages:          req.Languages,
		Certifications:     req.Certifications,
		Location:           req.Location,
		ResponseTime:      req.ResponseTime,
	})
	if err != nil {
		return nil, err
	}
	return &profile.UpdateProfileResponse{Message: "Profile updated successfully"}, nil
}

func (h *Handler) GetProfile(ctx context.Context, req *profile.GetProfileRequest) (*profile.GetProfileResponse, error) {
	if !h.roles.HasRole(ctx, "freelancer", "client", "admin") {
		return nil, ErrPermissionDenied
	}
	
	p, err := h.service.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &profile.GetProfileResponse{
		UserId:            p.UserID.String(),
		Title:              p.Title,
		Bio:                p.Bio,
		HourlyRate:        float32(p.HourlyRate),
		YearsOfExperience: int32(p.YearsOfExperience),
		Skills:             p.Skills,
		Languages:          p.Languages,
		Certifications:     p.Certifications,
		Location:           p.Location,
		ResponseTime:      p.ResponseTime,
	}, nil
}

