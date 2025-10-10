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
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{Client: rdb, Ctx: context.Background()}
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	return r.Client.Set(r.Ctx, key, value, ttl).Err()
}

func (r *RedisCache) Del(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}
