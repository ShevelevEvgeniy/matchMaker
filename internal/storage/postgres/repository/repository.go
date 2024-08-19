package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"matchMaker/internal/storage/postgres/repository/models"
	queryStr "matchMaker/internal/storage/postgres/repository/query"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Begin(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}

	return tx, nil
}

func (r *Repository) GetUsersInSearch(ctx context.Context, tx pgx.Tx, batchSize int) ([]models.User, error) {
	query := queryStr.GetUsersInSearch

	users := make([]models.User, 0, batchSize)
	rows, err := tx.Query(ctx, query, batchSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users in search")
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Skill, &user.Latency, &user.SearchStartTime); err != nil {
			fmt.Println(err)
			return nil, errors.Wrap(err, "failed to get users in search")
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to get users in search")
	}

	return users, nil
}

func (r *Repository) UnmarkSearch(ctx context.Context, tx pgx.Tx, usersIDs []int64) error {
	if len(usersIDs) == 0 {
		return nil
	}

	query := queryStr.UnmarkSearch

	_, err := tx.Exec(ctx, query, usersIDs)
	if err != nil {
		return errors.Wrap(err, "failed to unmark search")
	}

	return nil
}

func (r *Repository) SaveUsers(ctx context.Context, user models.User) error {
	query := queryStr.AddUser

	_, err := r.db.Exec(ctx, query, user.Name, user.Skill, user.Latency, user.SearchMatch, user.SearchStartTime)
	if err != nil {
		return errors.Wrap(err, "failed to add user")
	}

	return nil
}
