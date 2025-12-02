package configs

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	IPAddress string
	Password  string
	Database  int
}

type RedisInstance struct {
	RedisClient  *redis.Client
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
		RedisClient:  redisDatabase,
		redisContext: redisContext,
	}, nil
}

func (redisInstance *RedisInstance) Publish(topic string, message string) error {
	if redisInstance.RedisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	err := redisInstance.RedisClient.Publish(redisInstance.redisContext, topic, message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message to topic '%s': %w", topic, err)
	}

	logrus.Debug("📢 Published message to Redis topic '%s'\n", topic)
	return nil
}
