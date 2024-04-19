package interfaces

import (
	"context"
	"kiramishima/ionix/internal/models"
)

// AuthRepository interface
type AuthRepository interface {
	FindUserByCredentials(ctx context.Context, form *models.AuthForm) (*models.User, error)
	CreateAccount(ctx context.Context, form *models.RegisterForm) error
}
