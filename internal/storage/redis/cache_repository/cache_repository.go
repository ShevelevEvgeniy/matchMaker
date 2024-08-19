package cache_repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"matchMaker/config"
	"matchMaker/internal/storage/postgres/repository/models"
)

type CacheRepository struct {
	redis *redis.Client
	keys  config.Keys
}

func NewCacheRepository(redis *redis.Client, cfg config.Keys) *CacheRepository {
	return &CacheRepository{
		redis: redis,
		keys:  cfg,
	}
}

func (c *CacheRepository) GetRemainingUsers(ctx context.Context) ([]models.User, error) {
	usersBytes, err := c.redis.Get(ctx, c.keys.RemainingUsersKey).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get remaining users")
	}

	if len(usersBytes) == 0 {
		return nil, errors.Wrap(err, "not found remaining users")
	}

	var users []models.User
	err = json.Unmarshal(usersBytes, &users)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal users")
	}

	return users, nil
}

func (c *CacheRepository) SetRemainingUsers(ctx context.Context, users []models.User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return errors.Wrap(err, "failed to marshal users")
	}

	err = c.redis.Set(ctx, c.keys.RemainingUsersKey, data, time.Hour).Err()
	if err != nil {
		return errors.Wrap(err, "failed to set remaining users")
	}

	return nil
}

func (c *CacheRepository) DelRemainingUsers(ctx context.Context) error {
	err := c.redis.Del(ctx, c.keys.RemainingUsersKey).Err()
	if err != nil {
		return errors.Wrap(err, "failed to delete remaining users")
	}

	return nil
}
