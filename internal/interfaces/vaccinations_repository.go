package interfaces

import (
	"context"
	"kiramishima/ionix/internal/models"
)

// VaccinationRepository interface
type VaccinationRepository interface {
	GetVaccinationsData(ctx context.Context) ([]*models.Vaccination, error)
	CreateNewVaccinationItem(ctx context.Context, form *models.VaccinationForm) error
	GetVaccinationItemByID(ctx context.Context, vaccinationId int) (*models.Vaccination, error)
	UpdateVaccinationItem(ctx context.Context, vaccinationId int, form *models.Vaccination) error
	DeleteVaccinationItem(ctx context.Context, vaccinationId int) error
}
