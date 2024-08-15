package config

import "time"

type HTTPServer struct {
	Port        string        `envconfig:"HTTP_SERVER_PORT" env-required:"true"`
	Timeout     time.Duration `envconfig:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
	StopTimeout time.Duration `envconfig:"HTTP_SERVER_STOP_TIMEOUT" env-default:"5s"`
}
