package bootstrap

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"kiramishima/ionix/config"
	"kiramishima/ionix/internal/auth"
	"kiramishima/ionix/internal/drugs"
	"kiramishima/ionix/internal/pkg/database"
	"kiramishima/ionix/internal/server"
	"kiramishima/ionix/internal/vaccinations"
	"time"
)

func bootstrap(
	lifecycle fx.Lifecycle,
	logger *zap.SugaredLogger,
	server *server.Server,
) {

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("Starting API")
				return server.Run()
			},
			OnStop: func(ctx context.Context) error {
				return logger.Sync()
			},
		},
	)

}

var Module = fx.Options(
	config.Module,
	fx.Provide(func() *zap.SugaredLogger {
		logger, _ := zap.NewProduction()
		return logger.Sugar()
	}),
	fx.Provide(func() *chi.Mux {
		var r = chi.NewRouter()
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

		r.Use(middleware.Timeout(60 * time.Second))
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)
		r.Use(httprate.LimitByIP(1000, 1*time.Minute))
		r.Use(middleware.Compress(5))
		return r
	}),
	fx.Provide(func() *render.Render {
		return render.New()
	}),
	fx.Provide(func() *validator.Validate {
		return validator.New(validator.WithRequiredStructEnabled())
	}),
	server.Module,
	database.Module,
	auth.Module,
	drugs.Module,
	vaccinations.Module,
	fx.Invoke(bootstrap),
)
