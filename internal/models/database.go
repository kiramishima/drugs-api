package models

type Database struct {
	DatabaseDriver string `envconfig:"DATABASE_DRIVER" required:"true" default:"pgx"`
	DatabaseURL    string `envconfig:"DATABASE_URL" required:"true"`
	MaxOpenConns   int    `envconfig:"DATABASE_MAX_OPEN_CONNECTIONS" default:"25"`
	MaxIdleConns   int    `envconfig:"DATABASE_MAX_IDDLE_CONNECTIONS" default:"25"`
	MaxIdleTime    string `envconfig:"DATABASE_MAX_IDDLE_TIME" default:"15m"`
}
