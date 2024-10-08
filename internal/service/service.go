package service

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	userConverter "matchMaker/internal/converter/user_converter"
	"matchMaker/internal/dto"
	"matchMaker/internal/storage/postgres/repository/models"
)

type Repository interface {
	Begin(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	SaveUsers(ctx context.Context, user []models.User) error
	UnmarkSearch(ctx context.Context, tx pgx.Tx, usersId []int64) error
	GetUsersInSearch(ctx context.Context, tx pgx.Tx, groupSize int) ([]models.User, error)
}

type Cache interface {
	Watch(ctx context.Context, fn func(tx *redis.Tx) error) error
	SetRemainingUsers(ctx context.Context, users []models.User) error
	GetRemainingUsers(ctx context.Context, tx *redis.Tx) ([]models.User, error)
	DelRemainingUsers(ctx context.Context, tx *redis.Tx) error
	ExistsKey(ctx context.Context) bool
}
type Service struct {
	repo  Repository
	cache Cache
}

func NewService(repo Repository, cache Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) SaveUsers(ctx context.Context, users dto.Users) error {
	err := s.repo.SaveUsers(ctx, userConverter.ServiceToRepoModels(users))
	if err != nil {
		return errors.Wrap(err, "failed to save users")
	}

	return nil
}

func (s *Service) GetUsersInSearch(ctx context.Context, batchSize int) ([]models.User, error) {
	tx, err := s.repo.Begin(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}

	users, err := s.getUsersInSearchWithTx(ctx, tx, batchSize)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); err != nil {
			err = errors.Wrap(rollbackErr, "failed to rollback transaction")
		}

		return nil, errors.Wrap(err, "failed to get users in search")
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}

	return users, nil
}

func (s *Service) getUsersInSearchWithTx(ctx context.Context, tx pgx.Tx, batchSize int) ([]models.User, error) {
	users, err := s.repo.GetUsersInSearch(ctx, tx, batchSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users in search")
	}

	err = s.repo.UnmarkSearch(ctx, tx, userConverter.UsersToIds(users))
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmark search")
	}

	return users, nil
}

func (s *Service) SaveRemainingUsers(ctx context.Context, users []models.User) error {
	err := s.cache.SetRemainingUsers(ctx, users)
	if err != nil {
		return errors.Wrap(err, "failed to save remaining users")
	}

	return nil
}

func (s *Service) GetAndRemoveRemainingUsers(ctx context.Context) ([]models.User, bool, error) {
	var users []models.User

	ok := s.cache.ExistsKey(ctx)
	if !ok {
		return nil, false, nil
	}

	err := s.cache.Watch(ctx, func(tx *redis.Tx) error {
		var err error

		users, err = s.cache.GetRemainingUsers(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "failed to get remaining users")
		}

		err = s.cache.DelRemainingUsers(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "failed to delete remaining users")
		}

		return nil
	})

	if err != nil {
		return nil, false, errors.Wrap(err, "failed to watch cache")
	}

	return users, true, nil
}
