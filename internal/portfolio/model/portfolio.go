
package model

import "time"

type Portfolio struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FreelancerID  string    `gorm:"type:uuid;not null"`
	Title         string
	Description   string
	ImageURL      string
	Link          string
	CreatedAt     time.Time `gorm:"default:now()"`
}
