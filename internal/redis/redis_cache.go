package redis

import (
	"github.com/go-redis/redis/v8"
	"log"
	"context"
	"time"
)

var Rdb *redis.Client
var ctx = context.Background()

// Initialize connects to the Redis database
func Initialize() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Assuming Redis is running on localhost
	})
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
}

// SetOrderStatus caches the status of an order in Redis
func SetOrderStatus(orderID, status string) {
	err := Rdb.Set(ctx, orderID, status, 10*time.Minute).Err()
	if err != nil {
		log.Println("Error setting order status in Redis:", err)
	}
}

// GetOrderStatus retrieves the status of an order from Redis
func GetOrderStatus(orderID string) (string, error) {
	status, err := Rdb.Get(ctx, orderID).Result()
	if err == redis.Nil {
		return "", nil // Status not found in Redis
	} else if err != nil {
		log.Println("Error getting order status from Redis:", err)
		return "", err
	}
	return status, nil
}
