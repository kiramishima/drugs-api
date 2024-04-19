package drugs

import (
	"context"
	"errors"
	"go.uber.org/zap"
	impl "kiramishima/ionix/internal/interfaces"
	"kiramishima/ionix/internal/models"
	"time"
)

var _ impl.DrugService = (*service)(nil)

type service struct {
	logger         *zap.SugaredLogger
	repository     impl.DrugRepository
	contextTimeOut time.Duration
}

// NewDrugService creates a new auth service
func NewDrugService(repo impl.DrugRepository, logger *zap.SugaredLogger, timeout time.Duration) *service {
	return &service{
		logger:         logger,
		repository:     repo,
		contextTimeOut: timeout,
	}
}

func (svc service) GetListDrugs(ctx context.Context) ([]*models.Drug, error) {
	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()
	data, err := svc.repository.GetDrugsData(cxt)

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
				return nil, ErrServiceDrugs
			}
		}
	}
	// svc.logger.Info(data)
	return data, nil
}

func (svc service) NewDrug(ctx context.Context, form *models.DrugForm) error {
	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()

	err := svc.repository.CreateNewDrugItem(ctx, form)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-cxt.Done():
			return ErrTimeout
		default:
			if errors.Is(err, ErrDuplicateDrug) {
				return ErrDuplicateDrug
			} else {
				return ErrExecuteStatement
			}
		}
	}

	return nil
}

func (svc service) UpdateDrug(ctx context.Context, drugId int, form *models.DrugForm) error {
	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()
	// Retrieve the data
	drug, err := svc.repository.GetDrugItemByID(cxt, drugId)
	if errors.Is(err, ErrDrugNotFound) {
		return ErrDrugNotFound
	}
	svc.logger.Info(drug)
	//
	if form.Name != nil {
		drug.Name = *form.Name
	}
	if form.Approved != nil {
		drug.Approved = *form.Approved
	}
	if form.MinDose != nil {
		drug.MinDose = *form.MinDose
	}
	if form.MaxDose != nil {
		drug.MaxDose = *form.MaxDose
	}
	if form.AvailableAt != nil {
		var dt = *form.AvailableAt
		layout := "2006-01-02 15:04:05"
		tm, _ := time.Parse(layout, dt)
		svc.logger.Info(tm)
		drug.AvailableAt = tm
	}
	svc.logger.Info(form)
	// Call repository
	err = svc.repository.UpdateDrugItem(cxt, drugId, drug)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return ErrTimeout
		default:
			if errors.Is(err, ErrExecuteStatement) {
				return ErrExecuteStatement
			} else if errors.Is(err, ErrDrugNotFound) {
				return ErrDrugNotFound
			} else {
				return ErrUpdatingRecord
			}
		}
	}

	return nil
}

func (svc service) DeleteDrug(ctx context.Context, drugId int) error {
	cxt, cancel := context.WithTimeout(ctx, svc.contextTimeOut)
	defer cancel()

	_, err := svc.repository.GetDrugItemByID(cxt, drugId)
	if errors.Is(err, ErrDrugNotFound) {
		return ErrDrugNotFound
	}

	// Call repository
	err = svc.repository.DeleteDrugItem(cxt, drugId)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return ErrTimeout
		default:
			if errors.Is(err, ErrExecuteStatement) {
				return ErrExecuteStatement
			} else if errors.Is(err, ErrDrugNotFound) {
				return ErrDrugNotFound
			} else {
				return ErrDeletingRecord
			}
		}
	}

	return nil
}
