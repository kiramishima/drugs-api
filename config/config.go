package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"kiramishima/ionix/internal/models"
	"log"
)

// load loads the configuration from the environment
func load() (*models.Configuration, error) {
	var cfg models.Configuration
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// NewConfig creates and load config
func NewConfig() *models.Configuration {
	cfg, err := load()
	if err != nil {
		log.Printf("Can't load the configuration. Error: %s", err.Error())
	}

	return cfg
}

// Module config
var Module = fx.Options(
	fx.Provide(NewConfig),
)
