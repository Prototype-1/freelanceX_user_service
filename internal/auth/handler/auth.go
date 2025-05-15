package handler

import (
	"context"
	authPb "github.com/Prototype-1/freelanceX_user_service/proto/auth"
	"github.com/Prototype-1/freelanceX_user_service/internal/auth/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct {
	authPb.UnimplementedAuthServiceServer
	Service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

func (h *AuthHandler) Register(ctx context.Context, req *authPb.RegisterRequest) (*authPb.AuthResponse, error) {
	return h.Service.Register(ctx, req)
}

func (h *AuthHandler) Login(ctx context.Context, req *authPb.LoginRequest) (*authPb.AuthResponse, error) {
	return h.Service.Login(ctx, req)
}

func (h *AuthHandler) OAuthLogin(ctx context.Context, req *authPb.OAuthRequest) (*authPb.OAuthLoginResponse, error) {
	return h.Service.OAuthLogin(ctx, req)
}

func (h *AuthHandler) SelectRole(ctx context.Context, req *authPb.SelectRoleRequest) (*authPb.RoleSelectionResponse, error) {
	return h.Service.SelectRole(ctx, req)
}


func (h *AuthHandler) Logout(ctx context.Context, req *authPb.LogoutRequest) (*emptypb.Empty, error) {
	return h.Service.Logout(ctx, req)
}

func (h *AuthHandler) GetMe(ctx context.Context, req *authPb.SessionRequest) (*authPb.UserResponse, error) {
	return h.Service.GetMe(ctx, req)
}
