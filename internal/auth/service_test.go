package auth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/mocks"
	"kiramishima/ionix/internal/models"
	"testing"
	"time"
)

func TestService_SignIn(t *testing.T) {
	t.Parallel()
	zlog, _ := zap.NewProduction()
	logger := zlog.Sugar()
	mockCtrl := gomock.NewController(t)
	// ctx := context.Background()

	defer mockCtrl.Finish()

	repo := mocks.NewMockAuthRepository(mockCtrl)
	var good = &models.AuthForm{Email: "johnwick@gmail.com", Password: "123456"}
	tempGood := good.Hash256Password(good.Password)
	bcryptPass, _ := good.BcryptPassword(tempGood)

	var badPassword = &models.AuthForm{Email: "johnwick@gmail.com", Password: "23456"}

	var notExist = &models.AuthForm{Email: "kratos@gmail.com", Password: "123456"}

	var user = &models.User{
		ID:        1,
		Name:      "Jhon Wick",
		Email:     "jhonwick@gmail.com",
		Password:  bcryptPass,
		CreatedAt: time.Now(),
	}
	var user2 = &models.User{
		ID:        1,
		Name:      "Jhon Wick",
		Email:     "jhonwick@gmail.com",
		Password:  bcryptPass,
		CreatedAt: time.Now(),
	}
	// t.Log(good.ValidateBcryptPassword(user.Password, good.Password))

	repo.EXPECT().FindUserByCredentials(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
	repo.EXPECT().FindUserByCredentials(gomock.Any(), gomock.Any()).Times(1).Return(user2, ErrInvalidPassword)
	repo.EXPECT().FindUserByCredentials(gomock.Any(), notExist).Times(1).Return(nil, ErrUserNotFound)

	svc := NewAuthService(repo, logger, 5)

	t.Run("Good credentials", func(t *testing.T) {
		ctx := context.Background()
		var item, err = svc.SignIn(ctx, good)
		// t.Log(item, err)
		assert.NoError(t, err)
		assert.Equal(t, len(item.AccessToken) > 0, true)
	})

	t.Run("Bad credentials", func(t *testing.T) {
		ctx := context.Background()
		var item, err = svc.SignIn(ctx, badPassword)
		t.Log(item, err)
		assert.Error(t, err)
		// assert.Equal(t, len(item.AccessToken) > 0, true)
	})

	t.Run("User Not Found", func(t *testing.T) {
		ctx := context.Background()
		var item, err = svc.SignIn(ctx, notExist)
		t.Log(item, err)
		assert.Error(t, err)
		// assert.Equal(t, len(item.AccessToken) > 0, true)
	})
}
