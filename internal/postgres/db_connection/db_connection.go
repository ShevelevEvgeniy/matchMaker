package db_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"matchMaker/config"
)

func Connect(ctx context.Context, cfg config.DB) (*pgxpool.Pool, error) {
	urlExample := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		cfg.DriverName, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	connConfig, err := pgxpool.ParseConfig(urlExample)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse db url")
	}

	connConfig.MaxConns = cfg.MaxConns

	conCtx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	conn, err := pgxpool.NewWithConfig(conCtx, connConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to db")
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "Failed to ping db")
	}

	return conn, nil
}
