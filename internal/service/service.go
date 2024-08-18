package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"matchMaker/internal/converter"
	"matchMaker/internal/dto"
	"matchMaker/internal/storage/postgres/repository/models"
)

const (
	RemainingUsersKey = "remaining_users"
)

type Repository interface {
	Begin(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	Users(ctx context.Context, user models.User) error
	UnmarkSearch(ctx context.Context, usersId []int) error
	GetUsersInSearch(ctx context.Context, groupSize int) ([]models.User, error)
}

type Service struct {
	repo  Repository
	cache *redis.Client
}

func NewService(repo Repository, cache *redis.Client) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) Users(ctx context.Context, user dto.User) error {
	return s.repo.Users(ctx, converter.ServiceToRepoModel(user))
}

func (s *Service) GetUsersInSearch(ctx context.Context, batchSize int) ([]models.User, error) {
	tx, err := s.repo.Begin(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}

	ctx = context.WithValue(ctx, "tx", tx)
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Wrap(rollbackErr, "failed to rollback transaction")
			}
		}
	}()

	users, err := s.repo.GetUsersInSearch(ctx, batchSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users in search")
	}

	err = s.repo.UnmarkSearch(ctx, converter.UsersToIds(users))
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmark search")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}

	return users, nil
}

func (s *Service) SaveRemainingUsers(ctx context.Context, users []models.User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return errors.Wrap(err, "failed to marshal users")
	}

	err = s.cache.Set(ctx, RemainingUsersKey, data, time.Hour).Err()
	if err != nil {
		return errors.Wrap(err, "failed to save remaining users")
	}

	return nil
}
func (s *Service) GetAndRemoveRemainingUsers(ctx context.Context) ([]models.User, bool, error) {
	usersBytes, err := s.cache.Get(ctx, RemainingUsersKey).Bytes()
	if err != nil && err != redis.Nil {
		return nil, false, errors.Wrap(err, "failed to get remaining users")
	}

	if len(usersBytes) == 0 {
		return nil, false, nil
	}

	var users []models.User
	err = json.Unmarshal(usersBytes, &users)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to unmarshal users")
	}

	err = s.cache.Del(ctx, RemainingUsersKey).Err()
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to delete remaining users")
	}

	return nil, false, nil
}
