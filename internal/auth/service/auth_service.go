package service

import (
	"context"
	"errors"
	"time"

	authPb "github.com/Prototype-1/freelanceX_user_service/proto/auth"
	"github.com/Prototype-1/freelanceX_user_service/internal/auth/repository"
	"github.com/Prototype-1/freelanceX_user_service/pkg/jwt"
	"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
	"github.com/Prototype-1/freelanceX_user_service/internal/auth/model"
	"golang.org/x/crypto/bcrypt"
	//"log"
)

type AuthService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: repo}
}

func (s *AuthService) Register(ctx context.Context, req *authPb.RegisterRequest) (*authPb.AuthResponse, error) {
	if req.Role == "admin" {
		exists, err := s.UserRepo.IsAdminExists(ctx)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("An admin already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: new(string),
		Role:         req.Role,
	}

	*user.PasswordHash = string(hashedPassword)

	if err := s.UserRepo.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	sessionID := "session-id"
	if err := redis.SetSession(ctx, sessionID, user.ID.String(), time.Hour*24); err != nil {
		return nil, err
	}

	return &authPb.AuthResponse{
		AccessToken: accessToken,
		SessionId:   sessionID,
		UserId:      user.ID.String(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *authPb.LoginRequest) (*authPb.AuthResponse, error) {
	user, err := s.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	sessionID := "session-id"
	if err := redis.SetSession(ctx, sessionID, user.ID.String(), time.Hour*24); err != nil {
		return nil, err
	}

	return &authPb.AuthResponse{
		AccessToken: accessToken,
		SessionId:   sessionID,
		UserId:      user.ID.String(),
	}, nil
}

func (s *AuthService) OAuthLogin(ctx context.Context, req *authPb.OAuthRequest) (*authPb.OAuthLoginResponse, error) {
	user, err := s.UserRepo.GetUserByOAuthID(ctx, req.OauthProvider, req.OauthId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if !user.IsRoleSelected {
		return &authPb.OAuthLoginResponse{
			UserId: user.ID.String(),
		}, nil
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	sessionID := "session-id"
	if err := redis.SetSession(ctx, sessionID, user.ID.String(), time.Hour*24); err != nil {
		return nil, err
	}

	return &authPb.OAuthLoginResponse{
		Message:     "OAuth login successful",
		AccessToken: accessToken,
		SessionId:   sessionID,
		UserId:      user.ID.String(),
		IsRoleSelected: user.IsRoleSelected,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
	}, nil
}

func (s *AuthService) SelectRole(ctx context.Context, req *authPb.SelectRoleRequest) (*authPb.RoleSelectionResponse, error) {
	if req.Role != "freelancer" && req.Role != "client" && req.Role != "admin" {
		return nil, errors.New("invalid role")
	}

	if req.Role == "admin" {
		exists, err := s.UserRepo.IsAdminExists(ctx)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("An admin already exists")
		}
	}

	err := s.UserRepo.UpdateUserRole(ctx, req.UserId, req.Role)
	if err != nil {
		return nil, err
	}

	return &authPb.RoleSelectionResponse{
		Message: "Role updated successfully",
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *authPb.LogoutRequest) (*authPb.Empty, error) {
	if err := redis.DeleteSession(ctx, req.SessionId); err != nil {
		return nil, err
	}
	return &authPb.Empty{}, nil
}

func (s *AuthService) GetMe(ctx context.Context, req *authPb.SessionRequest) (*authPb.UserResponse, error) {
	claims, err := jwt.ParseAccessToken(req.Token)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return &authPb.UserResponse{
		Id:            user.ID.String(),
		Name:          user.Name,
		Email:         user.Email,
		Role:          user.Role,
		IsRoleSelected: user.IsRoleSelected,
	}, nil
}
