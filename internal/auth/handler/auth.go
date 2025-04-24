package handler

import (
	"context"
	"github.com/Prototype-1/freelanceX_user_service/internal/auth/service"
	"github.com/Prototype-1/freelanceX_user_service/pkg/jwt"
	"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	SessionService service.SessionService
}

func NewAuthHandler(sessionService service.SessionService) *AuthHandler {
	return &AuthHandler{SessionService: sessionService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	// Get email and password from the request
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Validate the credentials
	userID, err := validateUserCredentials(request.Email, request.Password) 
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 2. Generate JWT token (Access token)
	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		log.Println("Error generating JWT:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// 3. Create and save a session in Redis with a TTL
	sessionID := generateSessionID() // Create a unique session ID (UUID)
	ttl := time.Hour * 24 * 7 // Set TTL to 1 week for this session
	if err := redis.SetSession(c, sessionID, userID, ttl); err != nil {
		log.Println("Error saving session to Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// 4. Respond with the JWT and session ID
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"session_id":   sessionID,
	})
}

// Function to validate user credentials
func validateUserCredentials(email, password string) (string, error) {
	// Implement logic to fetch the user from the database by email
	// This is a dummy implementation; replace with actual logic
	// Assume we fetched user data from DB
	storedPasswordHash := "$2a$12$XKHq9T3mPz0i8aqME6nNeed.x6vV5W8lPf8.mHY4gf8Btz6LOl3dC" // Example hash for password "secret"
	userID := "user-uuid" // Replace with actual user ID from DB lookup

	// Compare password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
	if err != nil {
		return "", err // Invalid credentials
	}

	return userID, nil
}

// Function to generate a unique session ID
func generateSessionID() string {
	// Generate a UUID for the session ID
	sessionID := uuid.New().String() 
	return sessionID
}
