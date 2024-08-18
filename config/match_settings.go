package config

import "time"

type MatchSettings struct {
	GroupSize  int           `envconfig:"GROUP_SIZE" default:"5"`
	BatchSize  int           `envconfig:"BATCH" default:"100"`
	RangeSkill float64       `envconfig:"RANGE_SKILL" default:"1.0"`
	Delay      time.Duration `envconfig:"DELAY" default:"5s"`
}
