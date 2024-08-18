package config

type Redis struct {
	Protocol string `envconfig:"REDIS_PROTOCOL" default:"redis"`
	UserName string `envconfig:"REDIS_USER_NAME" required:"true"`
	Password string `envconfig:"REDIS_PASSWORD" required:"true"`
	Host     string `envconfig:"REDIS_HOST" required:"true"`
	Port     string `envconfig:"REDIS_PORT" required:"true"`
	DB       string `envconfig:"REDIS_DB" default:"0"`
	Params   string `envconfig:"REDIS_PARAMS" default:"?protocol=3"`
}
