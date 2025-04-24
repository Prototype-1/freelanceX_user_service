package service

import (
	"context"
	"time"
)

type SessionService interface {
	SetSession(ctx context.Context, sessionID, userID string, ttl time.Duration) error
	GetSession(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
