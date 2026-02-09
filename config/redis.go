package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	opt, err := redis.ParseURL(Config.RedisURL)
	if err != nil {
		log.Fatal("Invalid Redis URL:", err)
	}

	RedisClient = redis.NewClient(opt)

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Cannot connect to Redis:", err)
	}

	log.Println("Connected to Redis!")
}
