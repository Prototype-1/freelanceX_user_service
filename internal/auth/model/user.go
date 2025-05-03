package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(255)" json:"name"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   *string        `gorm:"type:varchar(255)" json:"-"`
	Role           string         `gorm:"type:varchar(20);not null" json:"role"` // 'admin', 'freelancer', 'client'
	OAuthProvider  *string        `gorm:"column:o_auth_provider;type:varchar(50)" json:"oauth_provider"`
	OAuthID        *string        `gorm:"column:o_auth_id;type:varchar(100)" json:"oauth_id"`
	IsRoleSelected bool           `gorm:"default:false" json:"is_role_selected"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // enables soft delete
}
