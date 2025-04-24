package model

import "time"

type FreelancerReview struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProjectID     string    `gorm:"type:uuid;not null;unique" json:"project_id"`
	FreelancerID  string    `gorm:"type:uuid;not null;index" json:"freelancer_id"`
	ClientID      string    `gorm:"type:uuid;not null" json:"client_id"`
	Rating        int       `gorm:"not null" json:"rating"`
	Feedback      string    `gorm:"type:text" json:"feedback"`
	CreatedAt     time.Time `gorm:"default:now()" json:"created_at"`
}
