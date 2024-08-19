package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	userConverter "matchMaker/internal/converter/user_converter"
	"matchMaker/internal/dto"
	"matchMaker/internal/storage/postgres/repository/models"
)

type Repository interface {
	Begin(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	SaveUsers(ctx context.Context, user models.User) error
	UnmarkSearch(ctx context.Context, tx pgx.Tx, usersId []int) error
	GetUsersInSearch(ctx context.Context, tx pgx.Tx, groupSize int) ([]models.User, error)
}

type Cache interface {
	SetRemainingUsers(ctx context.Context, users []models.User) error
	GetRemainingUsers(ctx context.Context) ([]models.User, error)
	DelRemainingUsers(ctx context.Context) error
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

func (s *Service) SaveUsers(ctx context.Context, user dto.User) error {
	return s.repo.SaveUsers(ctx, userConverter.ServiceToRepoModel(user))
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
	users, err := s.cache.GetRemainingUsers(ctx)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get remaining users")
	}

	err = s.cache.DelRemainingUsers(ctx)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to delete remaining users")
	}

	return users, false, nil
}
