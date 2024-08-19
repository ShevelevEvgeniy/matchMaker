package main

import (
	"context"
	"os"

	"go.uber.org/zap"
	"matchMaker/config"
	"matchMaker/internal/app"
	uberZap "matchMaker/lib/logger/uber_zap"
)

func main() {
	ctx := context.Background()

	log, stop, err := uberZap.InitLogger()
	if err != nil {
		os.Exit(1)
	}
	defer stop()

	cfg, err := config.MustLoad(log)
	if err != nil {
		log.Error("failed to load config", zap.Error(err))
		os.Exit(1)
	}

	log.Info("starting application")

	application := app.NewApp(log, cfg)
	if err = application.Run(ctx); err != nil {
		log.Error("failed to start app", zap.Error(err))
		os.Exit(1)
	}
}
