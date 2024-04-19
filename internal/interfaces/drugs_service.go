package interfaces

import (
	"context"
	models "kiramishima/ionix/internal/models"
)

// DrugService interface
type DrugService interface {
	GetListDrugs(ctx context.Context) ([]*models.Drug, error)
	NewDrug(ctx context.Context, form *models.DrugForm) error
	UpdateDrug(ctx context.Context, drugId int, form *models.DrugForm) error
	DeleteDrug(ctx context.Context, drugId int) error
}
