package uber_zap

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	EnvType string `envconfig:"ENV_TYPE" default:"development"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to load env config")
	}
	return cfg, nil
}
