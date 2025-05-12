package model

import (
	"time"
	"github.com/google/uuid"
	"github.com/lib/pq"	
)

type FreelancerProfile struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null"` 
	Title           string         `gorm:"type:varchar(255)"`
	Bio             string         `gorm:"type:text"`
	HourlyRate      float64        `gorm:"type:numeric"`
	YearsOfExperience int          `gorm:"type:int"`
Skills         pq.StringArray `gorm:"type:text[]"`
Languages      pq.StringArray `gorm:"type:text[]"`
Certifications pq.StringArray `gorm:"type:text[]"`
	Location        string         `gorm:"type:varchar(255)"`
	ResponseTime    string         `gorm:"type:varchar(50)"` 
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
}

