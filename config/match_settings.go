package config

import "time"

type MatchSettings struct {
	GroupSize    int           `envconfig:"GROUP_SIZE" default:"5"`
	BatchSize    int           `envconfig:"BATCH" default:"100"`
	CountWorkers int           `envconfig:"COUNT_WORKERS" default:"5"`
	Delay        time.Duration `envconfig:"DELAY" default:"5s"`
}
