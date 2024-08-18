package repository

import (
	"context"

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

func (r *Repository) getConn(ctx context.Context) *pgxpool.Pool {
	if conn, ok := ctx.Value("tx").(*pgxpool.Pool); ok {
		return conn
	}
	return r.db
}

func (r *Repository) Begin(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	return tx, nil
}

func (r *Repository) GetUsersInSearch(ctx context.Context, batchSize int) ([]models.User, error) {
	db := r.getConn(ctx)

	query := queryStr.GetUsersInSearch

	var users []models.User
	rows, err := db.Query(ctx, query, batchSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users in search")
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Skill, &user.Latency, &user.SearchMatch, &user.SearchStartTime); err != nil {
			return nil, errors.Wrap(err, "failed to get users in search")
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to get users in search")
	}

	return users, nil
}

func (r *Repository) UnmarkSearch(ctx context.Context, usersIDs []int) error {
	if len(usersIDs) == 0 {
		return nil
	}

	db := r.getConn(ctx)

	query := queryStr.UnmarkSearch

	params := make([]interface{}, 1)
	params[0] = usersIDs

	_, err := db.Exec(ctx, query, params...)
	if err != nil {
		return errors.Wrap(err, "failed to unmark search")
	}

	return nil
}

func (r *Repository) Users(ctx context.Context, user models.User) error {
	query := queryStr.AddUser

	_, err := r.db.Exec(ctx, query, user.Name, user.Skill, user.Latency, user.SearchMatch, user.SearchStartTime)
	if err != nil {
		return errors.Wrap(err, "failed to add user")
	}

	return nil
}
