package database

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/models"
	"time"
)

// NewDatabase creates an instance of DB
func NewDatabase(cfg *models.Configuration, logger *zap.SugaredLogger) (*sqlx.DB, error) {

	db, err := sqlx.Connect(cfg.DatabaseDriver, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	// conf connections
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	// iddle time
	duration, err := time.ParseDuration(cfg.MaxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ContextTimeout)*time.Second)
	defer cancel()

	// Ping to DB
	status := "up"
	err = db.PingContext(ctx)
	if err != nil {
		status = "down"
		return nil, err
	}
	logger.Debugf("Status DB: %s", status)
	return db, nil
}

var Module = fx.Module("db",
	fx.Provide(NewDatabase),
)
