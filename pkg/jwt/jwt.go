package jwt

import (
	"context"
	"errors"
	"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
	"log"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("your-jwt-secret-key")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateAccessToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiration
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

// Check if the session exists in Redis and matches the userID
func ValidateSession(sessionID, userID string) bool {
	ctx := context.Background()

	// Fetch session data from Redis
	storedUserID, err := redis.GetSession(ctx, sessionID)
	if err != nil {
		log.Println("Error fetching session from Redis:", err)
		return false
	}

	// If session is found, compare with userID
	if storedUserID != userID {
		log.Println("Session userID mismatch: expected", userID, "but got", storedUserID)
		return false
	}

	return true
}
