package redis_connection

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"matchMaker/config"
)

func Connect(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?protocol=%s",
		cfg.Protocol, cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DB, cfg.Params)

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse redis url")
	}
	client := redis.NewClient(opts)

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping redis server")
	}

	return client, nil
}
