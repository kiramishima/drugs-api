package interfaces

import (
	"context"
	models "kiramishima/ionix/internal/models"
)

// VaccinationService interface
type VaccinationService interface {
	GetListVaccinations(ctx context.Context) ([]*models.Vaccination, error)
	NewVaccination(ctx context.Context, form *models.VaccinationForm) error
	UpdateVaccination(ctx context.Context, vaccinationId int, form *models.VaccinationForm) error
	DeleteVaccination(ctx context.Context, vaccinationId int) error
}
