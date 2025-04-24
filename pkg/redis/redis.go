package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
	})
}

func GetRedisClient() *redis.Client {
	return rdb
}

func SetSession(ctx context.Context, sessionID string, userID string, ttl time.Duration) error {
	err := rdb.Set(ctx, sessionID, userID, ttl).Err()
	if err != nil {
		log.Println("Error setting session:", err)
		return err
	}
	return nil
}

func GetSession(ctx context.Context, sessionID string) (string, error) {
	userID, err := rdb.Get(ctx, sessionID).Result()
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
	err := rdb.Del(ctx, sessionID).Err()
	if err != nil {
		log.Println("Error deleting session:", err)
		return err
	}
	return nil
}
