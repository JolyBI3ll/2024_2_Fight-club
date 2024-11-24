package middleware

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
)

var RedisClient *redis.Client

func InitRedis() {
	address := os.Getenv("REDIS_ENDPOINT")
	password := os.Getenv("REDIS_PASSWORD")
	db := os.Getenv("REDIS_DB")

	dbNum := 0
	if db != "" {
		if num, err := strconv.Atoi(db); err == nil {
			dbNum = num
		} else {
			log.Fatalf("Invalid REDIS_DB value: %v", err)
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       dbNum,
	})

	ctx := context.Background()
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}
