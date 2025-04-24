package repository

import (
	"context"
	"errors"
	"github.com/Prototype-1/freelanceX_user_service/internal/auth/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetUserByOAuthID(ctx context.Context, oauthProvider, oauthID string) (*model.User, error)
	UpdateUserRole(ctx context.Context, userID string, role string) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByOAuthID(ctx context.Context, oauthProvider, oauthID string) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).Where("oauth_provider = ? AND oauth_id = ?", oauthProvider, oauthID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUserRole(ctx context.Context, userID string, role string) error {
	return r.DB.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"role":              role,
			"is_role_selected": true,
		}).Error
}

