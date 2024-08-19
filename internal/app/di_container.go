package app

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"matchMaker/config"
	"matchMaker/internal/http_server/api/v1/handlers"
	"matchMaker/internal/http_server/events"
	playerSelection "matchMaker/internal/http_server/handlers/player_selection"
	"matchMaker/internal/service"
	dbConn "matchMaker/internal/storage/postgres/db_connection"
	"matchMaker/internal/storage/postgres/repository"
	cacheRepository "matchMaker/internal/storage/redis/cache_repository"
	redisConn "matchMaker/internal/storage/redis/redis_connection"
)

type DIContainer struct {
	log                    *zap.Logger
	cfg                    *config.Config
	validator              *validator.Validate
	db                     *pgxpool.Pool
	redis                  *redis.Client
	cache                  service.Cache
	repo                   service.Repository
	serv                   handlers.Service
	handler                *handlers.MatchMakerHandler
	playerSelection        *playerSelection.PlayerSelection
	playerSelectionService playerSelection.Service
	formedGroupEvent       playerSelection.Events
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
	if di.redis == nil {
		redisClient, err := redisConn.Connect(ctx, di.cfg.Redis)
		if err != nil {
			di.log.Fatal("failed to connect to redis", zap.Error(err))
			os.Exit(1)
		}

		di.log.Info("connected to redis", zap.String("redis", di.cfg.Redis.Host))
		di.redis = redisClient
	}

	return di.redis
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

func (di *DIContainer) Cache(ctx context.Context) service.Cache {
	if di.cache == nil {
		di.cache = cacheRepository.NewCacheRepository(di.Redis(ctx), di.cfg.Redis.Keys)
	}

	return di.cache
}

func (di *DIContainer) Service(ctx context.Context) handlers.Service {
	if di.serv == nil {
		di.serv = service.NewService(di.Repository(ctx), di.Cache(ctx))
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
		di.playerSelectionService = service.NewService(di.Repository(ctx), di.Cache(ctx))
	}

	return di.playerSelectionService
}

func (di *DIContainer) FormedGroupEvent(_ context.Context) playerSelection.Events {
	if di.formedGroupEvent == nil {
		di.formedGroupEvent = events.NewFormedGroupEvent(di.log)
	}

	return di.formedGroupEvent
}

func (di *DIContainer) PlayerSelection(ctx context.Context) *playerSelection.PlayerSelection {
	if di.playerSelection == nil {
		di.playerSelection = playerSelection.NewPlayerSelection(di.cfg.MatchSettings, di.log, di.PlayerSelectionService(ctx), di.FormedGroupEvent(ctx))
	}

	return di.playerSelection
}
