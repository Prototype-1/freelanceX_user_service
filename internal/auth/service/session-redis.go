package service

import (
	"context"
	"github.com/Prototype-1/freelanceX_user_service/pkg/redis"
	"time"
)

type RedisSessionService struct{}

func NewRedisSessionService() *RedisSessionService {
	return &RedisSessionService{}
}

func (s *RedisSessionService) SetSession(ctx context.Context, sessionID, userID string, ttl time.Duration) error {
	return redisclient.SetSession(ctx, sessionID, userID, ttl)
}

func (s *RedisSessionService) GetSession(ctx context.Context, sessionID string) (string, error) {
	return redisclient.GetSession(ctx, sessionID)
}

func (s *RedisSessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return redisclient.DeleteSession(ctx, sessionID)
}
