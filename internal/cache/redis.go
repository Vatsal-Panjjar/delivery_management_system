package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) *RedisCache {
	r := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{
		client: r,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisCache) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
