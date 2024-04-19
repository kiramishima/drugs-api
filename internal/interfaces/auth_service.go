package interfaces

import (
	"context"
	models "kiramishima/ionix/internal/models"
)

// AuthService interface
type AuthService interface {
	SignIn(ctx context.Context, form *models.AuthForm) (*models.AuthResponse, error)
	SignUp(ctx context.Context, form *models.RegisterForm) error
}
