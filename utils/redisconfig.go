package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.TODO()

type RedisConfig struct{}

func (r *RedisConfig) InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping(ctx).Result()

	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Successfully connected to redis")
}
