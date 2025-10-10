package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisCache(addr string) *RedisCache {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisCache{Client: client, Ctx: ctx}
}

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	return r.Client.Set(r.Ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}
