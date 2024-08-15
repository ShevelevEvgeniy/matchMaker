package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"matchMaker/internal/postgres/repository/models"
	queryStr "matchMaker/internal/postgres/repository/query"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Users(ctx context.Context, user models.User) error {
	query := queryStr.AddUser

	_, err := r.db.Exec(ctx, query, user.Name, user.Skill, user.Latency, user.SearchingMatch, user.SearchStartTime)
	if err != nil {
		return errors.Wrap(err, "failed to add user")
	}

	return nil
}
