package vaccinations

import (
	"context"
	"errors"
	"go.uber.org/zap"
	impl "kiramishima/ionix/internal/interfaces"
	"kiramishima/ionix/internal/models"
	"time"
)

var _ impl.VaccinationService = (*service)(nil)

// NewVaccinationService creates a new vaccination service
func NewVaccinationService(repo impl.VaccinationRepository, logger *zap.Logger, timeout time.Duration) *service {
	return &service{
		logger:         logger,
		repository:     repo,
		contextTimeOut: timeout,
	}
}

type service struct {
	logger         *zap.Logger
	repository     impl.VaccinationRepository
	contextTimeOut time.Duration
}

func (svc service) GetListVaccinations(ctx context.Context) ([]*models.Vaccination, error) {
	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()
	data, err := svc.repository.GetVaccinationsData(cxt)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return nil, ErrTimeout
		default:
			if errors.Is(err, ErrNoRecords) {
				return nil, ErrNoRecords
			} else if errors.Is(err, ErrExecuteStatement) {
				return nil, ErrExecuteStatement
			} else {
				return nil, ErrServiceVaccination
			}
		}
	}
	svc.logger.Info("GetListVaccinations", zap.Any("data", data))
	return data, nil
}

func (svc service) NewVaccination(ctx context.Context, form *models.VaccinationForm) error {
	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()

	err := svc.repository.CreateNewVaccinationItem(ctx, form)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-cxt.Done():
			return ErrTimeout
		default:
			if errors.Is(err, ErrDuplicateVaccination) {
				return ErrDuplicateVaccination
			} else if errors.Is(err, ErrVaccinationNotFound) {
				return ErrVaccinationNotFound
			} else {
				return ErrExecuteStatement
			}
		}
	}

	return nil
}

func (svc service) UpdateVaccination(ctx context.Context, vaccinationId int, form *models.VaccinationForm) error {
	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()
	// Retrieve the data
	vaccination, err := svc.repository.GetVaccinationItemByID(cxt, vaccinationId)
	svc.logger.Info("UpdateVaccination", zap.Any("data", vaccination), zap.Any("err", err))

	if errors.Is(err, ErrVaccinationNotFound) {
		return ErrVaccinationNotFound
	}
	svc.logger.Info("UpdateVaccination", zap.Any("data", vaccination))

	//
	if form.Name != nil {
		vaccination.Name = *form.Name
	}
	if form.DrugID != nil {
		vaccination.DrugID = int32(*form.DrugID)
	}
	if form.Dose != nil {
		vaccination.Dose = int32(*form.Dose)
	}
	if form.AppliedAt != nil {
		var dt = *form.AppliedAt
		layout := "2006-01-02 15:04:05"
		tm, _ := time.Parse(layout, dt)
		svc.logger.Info("UpdateVaccination", zap.Any("tm", tm))
		vaccination.AppliedAt = tm
	}
	svc.logger.Info("UpdateVaccination", zap.Any("form", form))

	// Call repository
	err = svc.repository.UpdateVaccinationItem(cxt, vaccinationId, vaccination)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return ErrTimeout
		default:
			if errors.Is(err, ErrExecuteStatement) {
				return ErrExecuteStatement
			} else if errors.Is(err, ErrVaccinationNotFound) {
				return ErrVaccinationNotFound
			} else {
				return ErrUpdatingRecord
			}
		}
	}

	return nil
}

func (svc service) DeleteVaccination(ctx context.Context, vaccinationId int) error {
	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()

	_, err := svc.repository.GetVaccinationItemByID(cxt, vaccinationId)
	svc.logger.Info("DeleteVaccination", zap.Any("vaccinationId", vaccinationId), zap.Any("err", err))

	if errors.Is(err, ErrVaccinationNotFound) {
		return ErrVaccinationNotFound
	}

	// Call repository
	err = svc.repository.DeleteVaccinationItem(cxt, vaccinationId)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return ErrTimeout
		default:
			if errors.Is(err, ErrExecuteStatement) {
				return ErrExecuteStatement
			} else if errors.Is(err, ErrVaccinationNotFound) {
				return ErrVaccinationNotFound
			} else {
				return ErrDeletingRecord
			}
		}
	}

	return nil
}
