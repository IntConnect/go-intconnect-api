package configs

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	IPAddress string
	Password  string
	Database  int
}

type RedisInstance struct {
	redisClient  *redis.Client
	redisContext context.Context
}

func InitRedisInstance(redisConfig RedisConfig) (*RedisInstance, error) {
	redisContext := context.Background()
	redisDatabase := redis.NewClient(&redis.Options{
		Addr:     redisConfig.IPAddress,
		Password: redisConfig.Password,
		DB:       redisConfig.Database,
	})

	// Test connection
	if err := redisDatabase.Ping(redisContext).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect Redis: %w", err)
	}

	return &RedisInstance{
		redisClient:  redisDatabase,
		redisContext: redisContext,
	}, nil
}

func (redisInstance *RedisInstance) Publish(topic string, message string) error {
	if redisInstance.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	err := redisInstance.redisClient.Publish(redisInstance.redisContext, topic, message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message to topic '%s': %w", topic, err)
	}

	fmt.Printf("📢 Published message to Redis topic '%s'\n", topic)
	return nil
}
