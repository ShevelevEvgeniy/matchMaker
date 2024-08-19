package redis_connection

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"matchMaker/config"
)

func Connect(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	client := redis.NewClient(opts)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping redis server")
	}

	return client, nil
}
