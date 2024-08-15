package app

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"matchMaker/config"
	"matchMaker/internal/http_server/api/v1/handlers"
	dbConn "matchMaker/internal/postgres/db_connection"
	"matchMaker/internal/postgres/repository"
	"matchMaker/internal/service"
)

type DIContainer struct {
	log       *zap.Logger
	cfg       *config.Config
	validator *validator.Validate
	db        *pgxpool.Pool
	repo      service.Repository
	serv      handlers.Service
	handler   *handlers.MatchMakerHandler
}

func NewDiContainer(log *zap.Logger, cfg *config.Config) *DIContainer {
	return &DIContainer{
		log: log,
		cfg: cfg,
	}
}

func (di *DIContainer) DB(ctx context.Context) *pgxpool.Pool {
	if di.db == nil {
		db, err := dbConn.Connect(ctx, di.cfg.DB)
		if err != nil {
			di.log.Fatal("failed to connect to db", zap.Error(err))
			os.Exit(1)
		}

		di.log.Info("connected to db", zap.String("database", di.cfg.DB.DBName))
		di.db = db
	}

	return di.db
}

func (di *DIContainer) Validator() *validator.Validate {
	if di.validator == nil {
		di.validator = validator.New()
	}

	return di.validator
}
func (di *DIContainer) Repository(ctx context.Context) service.Repository {
	if di.repo == nil {
		di.repo = repository.NewRepository(di.DB(ctx))
	}

	return di.repo
}

func (di *DIContainer) Service(ctx context.Context) handlers.Service {
	if di.serv == nil {
		di.serv = service.NewService(di.Repository(ctx))
	}

	return di.serv
}

func (di *DIContainer) Handler(ctx context.Context) *handlers.MatchMakerHandler {
	if di.handler == nil {
		di.handler = handlers.NewMatchMakerHandler(di.log, di.Validator(), di.Service(ctx))
	}

	return di.handler

}
