package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Config struct {
	HTTPServer    HTTPServer
	DB            DB
	MatchSettings MatchSettings
	Redis         Redis
}

func MustLoad(log *zap.Logger) (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	logNonSecretConfig(log, &cfg)

	return &cfg, nil
}

func logNonSecretConfig(log *zap.Logger, cfg *Config) {
	log.Info("Initialized config",
		zap.String("http_server port", cfg.HTTPServer.Port),
		zap.Duration("http_server timeout", cfg.HTTPServer.Timeout),
		zap.Duration("http_server idle_timeout", cfg.HTTPServer.IdleTimeout),
		zap.Duration("http_server stop_timeout", cfg.HTTPServer.StopTimeout),
		zap.String("db host", cfg.DB.Host),
		zap.String("db port", cfg.DB.Port),
		zap.String("db name", cfg.DB.DBName),
		zap.String("db driverName", cfg.DB.DriverName),
		zap.String("redis host", cfg.Redis.Host),
		zap.String("redis port", cfg.Redis.Port),
	)
}
