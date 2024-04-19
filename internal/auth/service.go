package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	impl "kiramishima/ionix/internal/interfaces"
	"kiramishima/ionix/internal/models"
	"kiramishima/ionix/internal/pkg/utils"
	"time"
)

var _ impl.AuthService = (*service)(nil)

type service struct {
	logger         *zap.SugaredLogger
	repository     impl.AuthRepository
	contextTimeOut time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(repo impl.AuthRepository, logger *zap.SugaredLogger, timeout time.Duration) *service {
	return &service{
		logger:         logger,
		repository:     repo,
		contextTimeOut: timeout,
	}
}

func (svc service) SignIn(ctx context.Context, form *models.AuthForm) (*models.AuthResponse, error) {
	form.Password = form.Hash256Password(form.Password)

	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()

	user, err := svc.repository.FindUserByCredentials(ctx, form)
	svc.logger.Info(user, err)
	if err != nil {
		select {
		case <-cxt.Done():
			svc.logger.Info(errors.New("auth service is closed"), ctx.Err())
			return nil, ErrServiceAuth
		default:
			if errors.Is(err, ErrPrepapareQuery) {
				svc.logger.Info(ErrPrepapareQuery.Error())
				return nil, ErrPrepapareQuery
			} else if errors.Is(err, ErrUserNotFound) {
				svc.logger.Info(ErrUserNotFound.Error())
				return nil, ErrUserNotFound
			} else {
				svc.logger.Info(ErrServiceAuth.Error())
				return nil, ErrServiceAuth
			}
		}
	}

	// Check Password
	if !form.ValidateBcryptPassword(user.Password, form.Password) {
		svc.logger.Info(ErrInvalidPassword.Error())
		return nil, ErrInvalidPassword
	}

	// Generate Token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		svc.logger.Error(err.Error(), fmt.Sprintf("%T", err))
		return nil, jwt.ErrSignatureInvalid
	}

	return &models.AuthResponse{AccessToken: token}, nil
}

func (svc service) SignUp(ctx context.Context, form *models.RegisterForm) error {
	// Hash password
	form.Password = form.Hash256Password(form.Password)
	// Hash Bcrypt
	form.Password, _ = form.BcryptPassword(form.Password)

	// context
	ctx, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()
	// Call repository
	err := svc.repository.CreateAccount(ctx, form)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return ErrServiceAuth
		default:
			if errors.Is(err, ErrUserExist) {
				return ErrUserExist
			} else if errors.Is(err, ErrFailInsertUser) {
				return ErrFailInsertUser
			} else {
				return ErrServiceAuth
			}
		}
	}

	return nil
}
