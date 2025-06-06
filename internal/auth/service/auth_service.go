package service

import (
"context"
"errors"
 "time"
 "github.com/google/uuid"
"golang.org/x/crypto/bcrypt"
 "github.com/Prototype-1/freelanceX_user_service/internal/auth/model"
"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
authPb "github.com/Prototype-1/freelanceX_user_service/proto/auth"
 "github.com/Prototype-1/freelanceX_user_service/internal/auth/repository"
 "github.com/Prototype-1/freelanceX_user_service/pkg/oauth"
"github.com/Prototype-1/freelanceX_user_service/pkg/jwt"
"google.golang.org/protobuf/types/known/emptypb"
"fmt"
"log"
"google.golang.org/grpc/status"
"google.golang.org/grpc/codes"
"encoding/json"
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
			return nil, errors.New("an admin already exists")
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
		IsRoleSelected: true,
	}

	*user.PasswordHash = string(hashedPassword)
	if err := s.UserRepo.CreateUser(ctx, &user); err != nil {
		if err.Error() == "email already exists" {
			return nil, errors.New("this email is already registered")
		}
		return nil, err
	}

	return &authPb.AuthResponse{
		Message: "Registration successful. Please log in to continue....",
		UserId:  user.ID.String(),
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

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	sessionID := uuid.New().String()
	if err := redis.SetSession(ctx, sessionID, user.ID.String(), time.Hour*24); err != nil {
		return nil, fmt.Errorf("failed to set session in Redis: %w", err)
	}

	if err := redis.SetUserOnline(ctx, user.ID.String(), 5*time.Minute); err != nil {
    log.Printf("Failed to set user online: %v", err)
}

	return &authPb.AuthResponse{
		AccessToken: accessToken,
		SessionId:   sessionID,
		UserId:      user.ID.String(),
	}, nil
}

func (s *AuthService) OAuthLogin(ctx context.Context, req *authPb.OAuthRequest) (*authPb.OAuthLoginResponse, error) {
	if req.OauthProvider != "google" {
		return nil, errors.New("unsupported oauth provider")
	}

	token, err := oauth.GoogleConfig.Exchange(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	client := oauth.GoogleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	user, err := s.UserRepo.GetUserByOAuthID(ctx, req.OauthProvider, userInfo.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		existingUser, err := s.UserRepo.GetUserByEmail(ctx, userInfo.Email)
		if err == nil && existingUser != nil {
			return nil, errors.New("this email is already registered")
		}
		newUser := &model.User{
			ID:             uuid.New(),
			Email:          userInfo.Email,
			Name:           userInfo.Name,
			OAuthProvider:  &req.OauthProvider,
			OAuthID:        &userInfo.ID,
			IsRoleSelected: false,
		}

		if err := s.UserRepo.CreateUser(ctx, newUser); err != nil {
			return nil, err
		}

		return &authPb.OAuthLoginResponse{
			UserId:         newUser.ID.String(),
			IsRoleSelected: false,
			Message:        "User created. Role selection required.",
			Name:           newUser.Name,
			Email:          newUser.Email,
		}, nil
	}

	if !user.IsRoleSelected {
		return &authPb.OAuthLoginResponse{
			UserId:         user.ID.String(),
			IsRoleSelected: false,
			Message:        "Role selection required",
			Name:           user.Name,
			Email:          user.Email,
		}, nil
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	sessionID := uuid.New().String()
	if err := redis.SetSession(ctx, sessionID, user.ID.String(), 24*time.Hour); err != nil {
		return nil, fmt.Errorf("failed to set session in Redis: %w", err)
	}

	if err := redis.SetUserOnline(ctx, user.ID.String(), 5*time.Minute); err != nil {
    log.Printf("Failed to set user online: %v", err)
}

	return &authPb.OAuthLoginResponse{
		Message:         "OAuth login successful",
		AccessToken:     accessToken,
		SessionId:       sessionID,
		UserId:          user.ID.String(),
		IsRoleSelected:  true,
		Name:            user.Name,
		Email:           user.Email,
		Role:            user.Role,
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
			return nil, errors.New("an admin already exists")
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

func (s *AuthService) Logout(ctx context.Context, req *authPb.LogoutRequest) (*emptypb.Empty, error) {
	if err := redis.DeleteSession(ctx, req.SessionId); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthService) GetMe(ctx context.Context, req *authPb.SessionRequest) (*authPb.UserResponse, error) {
    userID := req.GetUserId()

    if userID == "" {
        return nil, errors.New("user_id is required")
    }

    user, err := s.UserRepo.GetUserByID(ctx, userID)
if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user not found")
	}

    return &authPb.UserResponse{
        Id:            user.ID.String(),
        Name:          user.Name,
        Email:         user.Email,
        Role:          user.Role,
        IsRoleSelected: user.IsRoleSelected,
    }, nil
}

func (s *AuthService) GetUserEmail(ctx context.Context, req *authPb.GetUserEmailRequest) (*authPb.GetUserEmailResponse, error) {
	user, err := s.UserRepo.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &authPb.GetUserEmailResponse{
		Email: user.Email,
	}, nil
}

