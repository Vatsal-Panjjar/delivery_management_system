package cache


import (
"context"
"github.com/go-redis/redis/v8"
"time"
)


type RedisCache struct{
Client *redis.Client
Ctx context.Context
}


func NewRedis(addr string) *RedisCache {
r := redis.NewClient(&redis.Options{Addr: addr})
return &RedisCache{Client: r, Ctx: context.Background()}
}


func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
return r.Client.Set(r.Ctx, key, value, ttl).Err()
}
func (r *RedisCache) Get(key string) (string, error) { return r.Client.Get(r.Ctx, key).Result() }
func (r *RedisCache) Del(key string) error { return r.Client.Del(r.Ctx, key).Err() }
