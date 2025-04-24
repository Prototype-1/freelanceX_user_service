package db

import (
	"fmt"
	"log"

	userModel "github.com/Prototype-1/freelanceX_user_service/internal/auth/model"
	profileModel "github.com/Prototype-1/freelanceX_user_service/internal/profile/model"
	portfolioModel "github.com/Prototype-1/freelanceX_user_service/internal/portfolio/model"
	reviewModel "github.com/Prototype-1/freelanceX_user_service/internal/review/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = AutoMigrate()
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&userModel.User{},
		&profileModel.FreelancerProfile{},
		&portfolioModel.Portfolio{},
		&reviewModel.FreelancerReview{},
	)
	if err != nil {
		log.Printf("AutoMigrate error: %v", err)
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database migration successful")
	return nil
}
