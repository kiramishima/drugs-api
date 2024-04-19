package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	var port = 3030
	var dbURL = "secret"
	os.Clearenv()
	os.Setenv("PORT", "3030")
	os.Setenv("DATABASE_URL", dbURL)

	t.Run("Load config", func(t *testing.T) {
		var cfg = NewConfig()
		assert.Equal(t, port, cfg.Port)
		assert.Equal(t, dbURL, cfg.DatabaseURL)
	})
}
