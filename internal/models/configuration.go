package models

// Configuration struct
type Configuration struct {
	HTTPServer
	Database
	ContextTimeout int `envconfig:"CONTEXT_TIMEOUT" default:"2"`
}
