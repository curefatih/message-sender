package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type RedisCache[T any] struct {
	client *redis.Client
	cfg    *viper.Viper
}

func NewRedisCache[T any](client *redis.Client, cfg *viper.Viper) *RedisCache[T] {
	return &RedisCache[T]{
		client: client,
		cfg:    cfg,
	}
}

var _ Cache[any] = &RedisCache[any]{}

// Get implements Cache.
func (r *RedisCache[T]) Get(ctx context.Context, key string) (*T, error) {
	var data T

	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		log.Error().Msgf("could not get all data from redis: %w", err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(result), &data); err != nil {
		log.Error().Msgf("unmarshal error from redis result: %w", err)
		return nil, err
	}

	return &data, nil
}

// Set implements Cache.
func (r *RedisCache[T]) Set(ctx context.Context, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Error().Msgf("could not marshal book to redis: %w", err)
		return err
	}

	set := r.client.Set(ctx, key, string(data), time.Duration(r.cfg.GetInt("cache.message.ttl_in_minutes"))*time.Minute)
	_, err = set.Result()

	return err
}

func SetupRedisClient(ctx context.Context, cfg *viper.Viper) *redis.Client {
	redisOpt := &redis.Options{
		Addr:     fmt.Sprint(cfg.Get("cache.redis.host")),
		Password: fmt.Sprint(cfg.Get("cache.redis.password")),
		DB:       0,
	}
	client := redis.NewClient(redisOpt)

	status := client.Ping(ctx)
	log.Info().Msgf("redis status: %q", status)

	return client
}
