package vaccinations

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"kiramishima/credit_assigner/internal/mocks"
	"kiramishima/credit_assigner/internal/models"
	"testing"
)

func TestCredit_Assign(t *testing.T) {
	zlog, _ := zap.NewProduction()
	logger := zlog.Sugar()
	mockCtrl := gomock.NewController(t)
	// ctx := context.Background()

	defer mockCtrl.Finish()

	repo := mocks.NewMockCreditRepository(mockCtrl)
	var credit3000 = &models.Credit{Invest: 3000, Credit300: 10, Credit500: 0, Credit700: 0, Status: 1}
	var credit6700 = &models.Credit{Invest: 6700, Credit300: 20, Credit500: 0, Credit700: 1, Status: 1}
	var credit9000 = &models.Credit{Invest: 9000, Credit300: 30, Credit500: 0, Credit700: 0, Status: 1}
	var credit400 = &models.Credit{Invest: 400, Credit300: 0, Credit500: 0, Credit700: 0, Status: 0}
	var credit50 = &models.Credit{Invest: 50, Credit300: 0, Credit500: 0, Credit700: 0, Status: 0}
	repo.EXPECT().RegisterAssign(gomock.Any(), credit3000).Return(nil)
	repo.EXPECT().RegisterAssign(gomock.Any(), credit6700).Return(nil)
	repo.EXPECT().RegisterAssign(gomock.Any(), credit9000).Return(nil)
	repo.EXPECT().RegisterAssign(gomock.Any(), credit400).Return(errors.New("investment error"))
	repo.EXPECT().RegisterAssign(gomock.Any(), credit50).Return(errors.New("investment needs be multiply of 100"))

	svc := NewCreditService(repo, logger, 2)

	t.Run("Test 3000", func(t *testing.T) {
		var invest int32 = 3000
		ctx := context.Background()
		var item, err = svc.Assign(ctx, invest)
		t.Log(item, err)
		assert.NoError(t, err)
		assert.Equal(t, item, credit3000)
	})

	t.Run("Test 6700", func(t *testing.T) {
		var invest int32 = 6700
		ctx := context.Background()
		var item, err = svc.Assign(ctx, invest)
		t.Log(item, err)
	})

	t.Run("Test 9000", func(t *testing.T) {
		var invest int32 = 9000
		ctx := context.Background()
		var item, err = svc.Assign(ctx, invest)
		t.Log(item, err)
	})

	t.Run("Test 400", func(t *testing.T) {
		var invest int32 = 400
		ctx := context.Background()
		var item, err = svc.Assign(ctx, invest)
		t.Log(item, err)
	})

	t.Run("Test 50", func(t *testing.T) {
		var invest int32 = 50
		ctx := context.Background()
		var item, err = svc.Assign(ctx, invest)
		t.Log(item, err)
		assert.Error(t, err)
	})
}

func TestService_Stats(t *testing.T) {
	zlog, _ := zap.NewProduction()
	logger := zlog.Sugar()
	mockCtrl := gomock.NewController(t)
	// ctx := context.Background()

	defer mockCtrl.Finish()

	repo := mocks.NewMockCreditRepository(mockCtrl)

	var item = &models.Stats{
		TotalAssigns:        26,
		TotalSuccessAssigns: 16,
		TotalFailAssigns:    10,
		AVGSuccessAssigns:   70.13,
		AVGFailAssigns:      29.87,
	}
	var item2 = &models.Stats{}
	repo.EXPECT().SumUp(gomock.Any()).Times(1).Return(item, nil)
	repo.EXPECT().SumUp(gomock.Any()).Times(1).Return(item2, nil)

	svc := NewCreditService(repo, logger, 2)

	t.Run("Get stats", func(t *testing.T) {
		ctx := context.Background()
		var obj, err = svc.Stats(ctx)
		t.Log(obj, err)
		assert.NoError(t, err)
		assert.Equal(t, item, obj)
	})

	t.Run("Get stats - No records", func(t *testing.T) {
		ctx := context.Background()
		var obj, err = svc.Stats(ctx)
		t.Log(obj, err)
		assert.NoError(t, err)
		assert.Equal(t, item2, obj)
	})
}
