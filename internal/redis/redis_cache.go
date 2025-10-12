package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var rdb *redis.Client
var ctx = context.Background()

// Initialize Redis connection
func Initialize() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

// SetOrderTracking stores order status in Redis cache
func SetOrderTracking(orderID string, status string) {
	rdb.Set(ctx, orderID, status, 0)
}

// GetOrderTracking retrieves order status from Redis cache
func GetOrderTracking(orderID string) string {
	status, err := rdb.Get(ctx, orderID).Result()
	if err == redis.Nil {
		return "No status found"
	} else if err != nil {
		log.Fatalf("Error getting order status: %v", err)
	}
	return status
}
