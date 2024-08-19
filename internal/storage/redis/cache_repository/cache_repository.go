package cache_repository

import (
	"context"
	"encoding/json"

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

func (c *CacheRepository) Watch(ctx context.Context, fn func(tx *redis.Tx) error) error {
	return c.redis.Watch(ctx, func(tx *redis.Tx) error {
		return fn(tx)
	}, c.keys.RemainingUsersKey)
}

func (c *CacheRepository) SetRemainingUsers(ctx context.Context, users []models.User) error {
	pipe := c.redis.Pipeline()

	for _, user := range users {
		data, err := json.Marshal(user)
		if err != nil {
			return errors.Wrap(err, "failed to marshal user")
		}

		err = pipe.RPush(ctx, c.keys.RemainingUsersKey, data).Err()
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to set remaining users")
	}

	return nil
}

func (c *CacheRepository) GetRemainingUsers(ctx context.Context, tx *redis.Tx) ([]models.User, error) {
	userBytesList, err := tx.LRange(ctx, c.keys.RemainingUsersKey, 0, -1).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get remaining users")
	}

	if len(userBytesList) == 0 {
		return nil, errors.Wrap(err, "not found remaining users")
	}

	var users []models.User
	for _, userBytes := range userBytesList {
		var user models.User
		err = json.Unmarshal([]byte(userBytes), &user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal user")
		}
		users = append(users, user)
	}

	return users, nil
}

func (c *CacheRepository) ExistsKey(ctx context.Context) bool {
	ok := c.redis.Exists(ctx, c.keys.RemainingUsersKey).Val()
	if ok != 1 {
		return false
	}

	return true
}

func (c *CacheRepository) DelRemainingUsers(ctx context.Context, tx *redis.Tx) error {
	err := tx.Del(ctx, c.keys.RemainingUsersKey).Err()
	if err != nil {
		return errors.Wrap(err, "failed to delete remaining users")
	}

	return nil
}
