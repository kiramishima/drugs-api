package vaccinations

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
var Module = fx.Module("vaccinations",
	fx.Invoke(func(conn *sqlx.DB, logger *zap.SugaredLogger, cfg *models.Configuration, r *chi.Mux, render *render.Render, validate *validator.Validate) error {
		// loads repository
		var repo = NewVaccinationRepository(conn, logger)
		// loads service
		var svc = NewVaccinationService(repo, logger, time.Duration(cfg.ContextTimeout)*time.Second)
		// loads handlers
		NewVaccionationHandlers(r, logger, svc, render, validate)
		return nil
	}),
)
