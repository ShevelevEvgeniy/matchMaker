package config

type Redis struct {
	Password string `envconfig:"REDIS_PASSWORD" required:"true"`
	Host     string `envconfig:"REDIS_HOST" required:"true"`
	Port     string `envconfig:"REDIS_PORT" required:"true"`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
	Keys     Keys
}

type Keys struct {
	RemainingUsersKey string `envconfig:"REMAINING_USERS_KEY" required:"true"`
}
