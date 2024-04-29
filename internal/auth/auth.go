package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/models"
	"time"
)

// Module auth
var Module = fx.Module("auth",
	fx.Invoke(func(conn *sqlx.DB, logger *zap.Logger, cfg *models.Configuration, r *chi.Mux, render *render.Render, validate *validator.Validate) error {
		// loads repository
		var repo = NewAuthRepository(conn, logger)
		// loads service
		var svc = NewAuthService(repo, logger, time.Duration(cfg.ContextTimeout)*time.Second)
		// loads handlers
		NewAuthHandlers(r, logger, svc, render, validate)
		return nil
	}),
)
