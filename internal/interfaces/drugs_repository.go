package interfaces

import (
	"context"
	"kiramishima/ionix/internal/models"
)

// DrugRepository interface
type DrugRepository interface {
	GetDrugsData(ctx context.Context) ([]*models.Drug, error)
	CreateNewDrugItem(ctx context.Context, form *models.DrugForm) error
	GetDrugItemByID(ctx context.Context, drugId int) (*models.Drug, error)
	UpdateDrugItem(ctx context.Context, drugId int, form *models.Drug) error
	DeleteDrugItem(ctx context.Context, drugId int) error
}
