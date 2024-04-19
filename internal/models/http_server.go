package models

import "time"

type HTTPServer struct {
	ServerAddress string        `envconfig:"SERVE_ADDRESS" default:"0.0.0.0"`
	IdleTimeout   time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT" default:"60s"`
	Port          int           `envconfig:"PORT" default:"8080"`
	ReadTimeout   time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"1s"`
	WriteTimeout  time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"2s"`
}
