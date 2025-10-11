package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	C   *redis.Client
	Ctx context.Context
}

// NewRedis creates a new Redis client connected to the given address.
func NewRedis(addr string) *RedisClient {
	c := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisClient{C: c, Ctx: context.Background()}
}

// Set stores a value in Redis with a TTL.
func (r *RedisClient) Set(key string, val interface{}, ttl time.Duration) error {
	return r.C.Set(r.Ctx, key, val, ttl).Err()
}

// Get retrieves a value from Redis.
func (r *RedisClient) Get(key string) (string, error) {
	return r.C.Get(r.Ctx, key).Result()
}

// Del deletes a key from Redis.
func (r *RedisClient) Del(key string) error {
	return r.C.Del(r.Ctx, key).Err()
}
