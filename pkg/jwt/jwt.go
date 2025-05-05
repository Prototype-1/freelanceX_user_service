package jwt

import (
	"context"
	"errors"
	"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
	"github.com/Prototype-1/freelanceX_user_service/config"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

func getSecretKey() []byte {
	return []byte(config.AppConfig.JWTSecret)
}

type Claims struct {
	UserID string `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), 
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

func ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if ve, ok := err.(jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, errors.New("token expired")
		}
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func ValidateSession(sessionID, userID string) bool {
	ctx := context.Background()

	storedUserID, err := redis.GetSession(ctx, sessionID)
	if err != nil {
		log.Println("Error fetching session from Redis:", err)
		return false
	}

	if storedUserID != userID {
		log.Println("Session userID mismatch: expected", userID, "but got", storedUserID)
		return false
	}

	return true
}
