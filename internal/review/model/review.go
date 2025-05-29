package model

import (
	"time"
	"github.com/google/uuid"
)

type FreelancerReview struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProjectID     uuid.UUID `gorm:"type:uuid;not null;unique" json:"project_id"`
	FreelancerID  uuid.UUID `gorm:"type:uuid;not null;index" json:"freelancer_id"`
	ClientID      uuid.UUID `gorm:"type:uuid;not null" json:"client_id"`
	Rating        int       `gorm:"not null" json:"rating"`
	Feedback      string    `gorm:"type:text" json:"feedback"`
	CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`
}
