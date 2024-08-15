package config

type DB struct {
	Host       string `envconfig:"DB_HOST" env-required:"true"`
	Port       string `envconfig:"DB_PORT" env-required:"true"`
	DBName     string `envconfig:"DB_NAME" env-required:"true"`
	Username   string `envconfig:"DB_USER_NAME" env-required:"true"`
	Password   string `envconfig:"DB_PASSWORD" env-required:"true"`
	SslMode    string `envconfig:"DB_SSL_MODE" env-default:"disable"`
	DriverName string `envconfig:"DB_DRIVER_NAME" env-default:"postgres"`
	MaxConns   int32  `envconfig:"DB_MAX_CONNS" env-default:"10"`
}
