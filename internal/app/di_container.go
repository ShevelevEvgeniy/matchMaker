package app

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"matchMaker/config"
	"matchMaker/internal/console"
	"matchMaker/internal/http_server/api/v1/handlers"
	playerSelection "matchMaker/internal/http_server/player_selection"
	"matchMaker/internal/service"
	dbConn "matchMaker/internal/storage/postgres/db_connection"
	"matchMaker/internal/storage/postgres/repository"
	"matchMaker/internal/storage/redis/redis_client"
)

type DIContainer struct {
	log                    *zap.Logger
	cfg                    *config.Config
	validator              *validator.Validate
	db                     *pgxpool.Pool
	cache                  *redis.Client
	repo                   service.Repository
	serv                   handlers.Service
	handler                *handlers.MatchMakerHandler
	playerSelection        playerSelection.PlayerSelectionInterface
	playerSelectionService playerSelection.Service
	logGroups              playerSelection.LogGroups
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

func (di *DIContainer) Redis(ctx context.Context) *redis.Client {
	if di.cache == nil {
		cache, err := redis_client.Connect(ctx, di.cfg.Redis)
		if err != nil {
			di.log.Fatal("failed to connect to redis", zap.Error(err))
			os.Exit(1)
		}

		di.log.Info("connected to redis", zap.String("redis", di.cfg.Redis.Host))
		di.cache = cache
	}

	return di.cache
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
		di.serv = service.NewService(di.Repository(ctx), di.Redis(ctx))
	}

	return di.serv
}

func (di *DIContainer) Handler(ctx context.Context) *handlers.MatchMakerHandler {
	if di.handler == nil {
		di.handler = handlers.NewMatchMakerHandler(di.log, di.Validator(), di.Service(ctx))
	}

	return di.handler

}

func (di *DIContainer) PlayerSelectionService(ctx context.Context) playerSelection.Service {
	if di.playerSelectionService == nil {
		di.playerSelectionService = service.NewService(di.Repository(ctx), di.Redis(ctx))
	}

	return di.playerSelectionService
}

func (di *DIContainer) LogGroups(_ context.Context) playerSelection.LogGroups {
	if di.logGroups == nil {
		di.logGroups = console.NewLogGroups(di.log)
	}

	return di.logGroups
}

func (di *DIContainer) PlayerSelection(ctx context.Context) playerSelection.PlayerSelectionInterface {
	if di.playerSelection == nil {
		di.playerSelection = playerSelection.NewPlayerSelection(di.cfg.MatchSettings, di.log, di.PlayerSelectionService(ctx), di.LogGroups(ctx))
	}

	return di.playerSelection
}
