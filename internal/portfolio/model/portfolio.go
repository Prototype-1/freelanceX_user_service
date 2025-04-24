package model

import "time"

type Portfolio struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	FreelancerID string    `gorm:"type:uuid;not null;index" json:"freelancer_id"`

	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	ImageURL    string `gorm:"type:text" json:"image_url"`
	Link        string `gorm:"type:text" json:"link"`

	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
}
