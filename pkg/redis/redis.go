package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
	"fmt"
	"github.com/Prototype-1/freelanceX_user_service/config"
)

var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: config.AppConfig.RedisAddr,
	})
	    ctx := context.Background()
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Println("Connected to Redis at:", config.AppConfig.RedisAddr)
}

func GetRedisClient() *redis.Client {
	return rdb
}

func SetSession(ctx context.Context, sessionID string, userID string, ttl time.Duration) error {
	key := "session:" + sessionID
	err := rdb.Set(ctx, key, userID, ttl).Err()
	if err != nil {
		log.Println("Error setting session:", err)
		return err
	}
	return nil
}

func GetSession(ctx context.Context, sessionID string) (string, error) {
	key := "session:" + sessionID
	userID, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil 
	}
	if err != nil {
		log.Println("Error getting session:", err)
		return "", err
	}
	return userID, nil
}

func DeleteSession(ctx context.Context, sessionID string) error {
	key := "session:" + sessionID
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Println("Error deleting session:", err)
		return err
	}
	return nil
}

func SetUserOnline(ctx context.Context, userID string, ttl time.Duration) error {
	key := fmt.Sprintf("online:%s", userID)
	return rdb.Set(ctx, key, "1", ttl).Err()
}