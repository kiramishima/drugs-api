package vaccinations

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

	repo := mocks.NewMockVaccinationRepository(mockCtrl)

	var data = []*models.Vaccination{
		{
			ID:        1,
			Name:      "Jhon Wick",
			Drug:      "medicament 1",
			DrugID:    1,
			Dose:      1,
			AppliedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Jhon Connor",
			Drug:      "medicament 1",
			DrugID:    1,
			Dose:      1,
			AppliedAt: time.Now(),
		},
	}

	repo.EXPECT().GetVaccinationsData(gomock.Any()).Times(1).Return(data, nil)
	repo.EXPECT().GetVaccinationsData(gomock.Any()).Times(1).Return(nil, ErrNoRecords)

	svc := NewVaccinationService(repo, logger, 5)

	t.Run("Ok- Getting Data", func(t *testing.T) {
		ctx := context.Background()
		var item, err = svc.GetListVaccinations(ctx)
		t.Log(item, err)
		assert.NoError(t, err)
		assert.Equal(t, len(item) > 0, true)
	})

	t.Run("Ok - No Rows", func(t *testing.T) {
		ctx := context.Background()
		var item, err = svc.GetListVaccinations(ctx)
		t.Log(item, err)
		assert.Error(t, err)
		assert.Equal(t, len(item) == 0, true)
	})
}
