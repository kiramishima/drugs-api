package drugs

import (
	"context"
	"kiramishima/ionix/internal/mocks"
	"kiramishima/ionix/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestService_GetListDrugs(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	// ctx := context.Background()

	defer mockCtrl.Finish()

	repo := mocks.NewMockDrugRepository(mockCtrl)

	var drugs = []*models.Drug{
		{
			ID:          1,
			Name:        "medicament 1",
			Approved:    true,
			MinDose:     1,
			MaxDose:     5,
			AvailableAt: time.Now(),
		},
		{
			ID:          2,
			Name:        "medicament 2",
			Approved:    true,
			MinDose:     1,
			MaxDose:     5,
			AvailableAt: time.Now(),
		},
	}
	// t.Log(good.ValidateBcryptPassword(user.Password, good.Password))

	repo.EXPECT().GetDrugsData(gomock.Any()).Times(1).Return(drugs, nil)
	repo.EXPECT().GetDrugsData(gomock.Any()).Times(1).Return(nil, ErrNoRecords)
	//repo.EXPECT().FindUserByCredentials(gomock.Any(), notExist).Times(1).Return(nil, ErrUserNotFound)

	svc := NewDrugService(repo, logger, 5)

	t.Run("Ok- Getting Data", func(t *testing.T) {
		ctx := context.Background()
		var item, err = svc.GetListDrugs(ctx)
		t.Log(item, err)
		assert.NoError(t, err)
		assert.Equal(t, len(item) > 0, true)
	})

	t.Run("Ok - No Rows", func(t *testing.T) {
		ctx := context.Background()
		var item, err = svc.GetListDrugs(ctx)
		t.Log(item, err)
		assert.Error(t, err)
		assert.Equal(t, len(item) == 0, true)
	})
}

func TestService_NewDrug(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	// ctx := context.Background()

	defer mockCtrl.Finish()

	repo := mocks.NewMockDrugRepository(mockCtrl)

	var name = "Aspirina"
	var approved = true
	var minDose = 1
	var maxDose = 2
	var availableAt = time.Now().String()
	var item = &models.DrugForm{
		Name:        &name,
		Approved:    &approved,
		MinDose:     &minDose,
		MaxDose:     &maxDose,
		AvailableAt: &availableAt,
	}
	// t.Log(good.ValidateBcryptPassword(user.Password, good.Password))

	repo.EXPECT().CreateNewDrugItem(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	repo.EXPECT().CreateNewDrugItem(gomock.Any(), gomock.Any()).Times(1).Return(ErrDuplicateDrug)
	//repo.EXPECT().FindUserByCredentials(gomock.Any(), notExist).Times(1).Return(nil, ErrUserNotFound)

	svc := NewDrugService(repo, logger, 5)

	t.Run("Ok- Getting Data", func(t *testing.T) {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
		var err = svc.NewDrug(ctx, item)
		t.Log(err)
		assert.NoError(t, err)
		// assert.Equal(t, len(item) > 0, true)
	})

	t.Run("Duplicate record", func(t *testing.T) {
		ctx := context.Background()
		var err = svc.NewDrug(ctx, item)
		t.Log(item, err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrDuplicateDrug.Error())
	})
}

func TestService_UpdateDrug(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	// ctx := context.Background()

	defer mockCtrl.Finish()

	repo := mocks.NewMockDrugRepository(mockCtrl)

	var name = "Aspirina"
	var approved = true
	var minDose = 1
	var maxDose = 2
	// var availableAt = time.Now().String()
	var item = &models.DrugForm{
		Name:        &name,
		Approved:    &approved,
		MinDose:     &minDose,
		MaxDose:     &maxDose,
		AvailableAt: nil,
	}
	// Record
	var drug = &models.Drug{
		ID:          1,
		Name:        "medicament 1",
		Approved:    true,
		MinDose:     1,
		MaxDose:     5,
		AvailableAt: time.Now(),
	}

	svc := NewDrugService(repo, logger, 5)

	t.Run("Ok- Updating Data", func(t *testing.T) {
		ctx := context.Background()
		id := 1

		repo.EXPECT().GetDrugItemByID(gomock.Any(), id).Times(1).Return(drug, nil)
		repo.EXPECT().UpdateDrugItem(gomock.Any(), id, gomock.Any()).Times(1).Return(nil)

		var err = svc.UpdateDrug(ctx, id, item)
		t.Log(err)
		assert.NoError(t, err)
		// assert.Equal(t, len(item) > 0, true)
	})

	t.Run("No existing record", func(t *testing.T) {
		ctx := context.Background()

		id := 2
		repo.EXPECT().GetDrugItemByID(gomock.Any(), gomock.Any()).Times(1).Return(nil, ErrDrugNotFound)

		var err = svc.UpdateDrug(ctx, id, item)
		t.Log(err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrDrugNotFound.Error())
	})
}

func TestService_DeleteDrug(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	// ctx := context.Background()

	defer mockCtrl.Finish()

	repo := mocks.NewMockDrugRepository(mockCtrl)

	/*var name = "Aspirina"
	var approved = true
	var minDose = 1
	var maxDose = 2
	// var availableAt = time.Now().String()
	/*var item = &models.DrugForm{
		Name:        &name,
		Approved:    &approved,
		MinDose:     &minDose,
		MaxDose:     &maxDose,
		AvailableAt: nil,
	}*/
	// Record
	var drug = &models.Drug{
		ID:          1,
		Name:        "medicament 1",
		Approved:    true,
		MinDose:     1,
		MaxDose:     5,
		AvailableAt: time.Now(),
	}

	// Drug

	// Update

	// repo.EXPECT().DeleteDrugItem(gomock.Any(), gomock.Eq(1)).Times(1).Return(nil)
	//repo.EXPECT().FindUserByCredentials(gomock.Any(), notExist).Times(1).Return(nil, ErrUserNotFound)

	svc := NewDrugService(repo, logger, 5)

	t.Run("Ok- Deleting Data", func(t *testing.T) {
		ctx := context.Background()

		id := 1

		repo.EXPECT().GetDrugItemByID(gomock.Any(), id).Times(1).Return(drug, nil)
		repo.EXPECT().DeleteDrugItem(gomock.Any(), id).Times(1).Return(nil)
		var err = svc.DeleteDrug(ctx, id)
		t.Log(err)
		assert.NoError(t, err)
		// assert.Equal(t, len(item) > 0, true)
	})

	t.Run("No existing record", func(t *testing.T) {
		ctx := context.Background()
		id := 1
		repo.EXPECT().GetDrugItemByID(gomock.Any(), gomock.Any()).Times(1).Return(nil, ErrDrugNotFound)

		var err = svc.DeleteDrug(ctx, id)
		t.Log(err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrDrugNotFound.Error())
	})
}
