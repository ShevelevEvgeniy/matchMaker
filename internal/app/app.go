package app

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"matchMaker/config"
)

type App struct {
	log *zap.Logger
	cfg *config.Config
}

func NewApp(log *zap.Logger, cfg *config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	di := NewDiContainer(a.log, a.cfg)

	router := initRouter(ctx, di)

	server := NewServer(a.cfg, router)
	err := server.Run(a.log, a.cfg)
	if err != nil {
		a.log.Error("error occurred on http_server shutting down:", zap.String("error", err.Error()))
		return errors.Wrap(err, "error occurred on http_server shutting down")
	}

	a.log.Info("application started")

	server.Shutdown(ctx, a.log, a.cfg.HTTPServer.StopTimeout)

	a.log.Info("application stopped")

	return nil
}
