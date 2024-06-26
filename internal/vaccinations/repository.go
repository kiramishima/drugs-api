package vaccinations

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/interfaces"
	"kiramishima/ionix/internal/models"
)

// implement drug repository
var _ interfaces.VaccinationRepository = (*repository)(nil)

// NewVaccinationRepository Creates a new instance of Repository
func NewVaccinationRepository(conn *sqlx.DB, logger *zap.Logger) *repository {
	return &repository{
		db:  conn,
		log: logger,
	}
}

// Repository struct
type repository struct {
	db  *sqlx.DB
	log *zap.Logger
}

func (repo repository) GetVaccinationsData(ctx context.Context) ([]*models.Vaccination, error) {
	var query = `SELECT
		v.id,
		v.name,
		d.name drug,
		v.drug_id,
		v.dose,
		v.applied_at
	FROM vaccinations v
	INNER JOIN drugs d on d.id = v.drug_id
	WHERE d.deleted_at IS NULL OR v.deleted_at IS NULL`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, ErrPrepapareQuery
	}
	defer func(stmt *sqlx.Stmt) {
		err := stmt.Close()
		if err != nil {
			repo.log.Error("failed to close statement", zap.Error(err))
		}
	}(stmt)

	var list = make([]*models.Vaccination, 0)

	rows, err := stmt.QueryxContext(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return list, nil
		} else {
			return list, ErrExecuteStatement
		}
	}
	for rows.Next() {
		var appliedAt sql.NullTime
		var item = &models.Vaccination{}
		err = rows.Scan(&item.ID, &item.Name, &item.Drug, &item.DrugID, &item.Dose, &appliedAt)

		if errors.Is(err, sql.ErrNoRows) {
			break
		}

		if appliedAt.Valid {
			item.AppliedAt = appliedAt.Time
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo repository) CreateNewVaccinationItem(ctx context.Context, form *models.VaccinationForm) error {
	repo.log.Info("[INFO]", zap.Any("form", form))
	var query = `INSERT INTO vaccinations (name, drug_id, dose, applied_at)
	VALUES ($1, $2, $3, CAST($4 AS TIMESTAMP))`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return ErrPrepapareQuery
	}
	defer func(stmt *sqlx.Stmt) {
		err := stmt.Close()
		if err != nil {
			repo.log.Error("failed to close statement", zap.Error(err))
		}
	}(stmt)

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return ErrBeginTransaction
	}
	defer func(tx *sqlx.Tx) {
		err := tx.Rollback()
		if err != nil {
			repo.log.Error("failed to rollback", zap.Error(err))
		}
	}(tx)

	_, err = stmt.ExecContext(ctx, form.Name, form.DrugID, form.Dose, form.AppliedAt)

	if err != nil {
		repo.log.Info(err.Error())
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok {
			repo.log.Info(string(pgErr.Code))
			if pgErr.Code == "23505" {
				return ErrDuplicateVaccination
			} else {
				return ErrInsertFailed
			}
		}

	}
	if err = tx.Commit(); err != nil {
		return ErrCommitTransaction
	}
	return nil
}

func (repo repository) GetVaccinationItemByID(ctx context.Context, vaccinationId int) (*models.Vaccination, error) {
	var query = `SELECT
		v.id,
		v.name,
		d.name drug,
		v.drug_id,
		v.dose,
		v.applied_at
	FROM vaccinations v
	INNER JOIN drugs d on d.id = v.drug_id
	WHERE (d.deleted_at IS NULL OR v.deleted_at IS NULL) AND v.id = $1`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, ErrPrepapareQuery
	}
	defer func(stmt *sqlx.Stmt) {
		err := stmt.Close()
		if err != nil {
			repo.log.Error("failed to prepare statement", zap.Error(err))
		}
	}(stmt)

	row := stmt.QueryRowContext(ctx, vaccinationId)

	var appliedAt sql.NullTime
	var item = &models.Vaccination{}
	err = row.Scan(&item.ID, &item.Name, &item.Drug, &item.DrugID, &item.Dose, &appliedAt)
	repo.log.Info("[INFO]", zap.Any("item", item))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrVaccinationNotFound
	}

	if appliedAt.Valid {
		item.AppliedAt = appliedAt.Time
	}

	return item, nil
}

func (repo repository) UpdateVaccinationItem(ctx context.Context, vaccinationId int, form *models.Vaccination) error {
	repo.log.Info("[INFO]", zap.Any("form", form))
	var query = `UPDATE vaccinations SET name = $1, drug_id = $2, dose = $3, applied_at = $4, updated_at=NOW() WHERE id = $5`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return ErrPrepapareQuery
	}
	defer func(stmt *sqlx.Stmt) {
		err := stmt.Close()
		if err != nil {
			repo.log.Error("failed to prepare statement", zap.Error(err))
		}
	}(stmt)

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return ErrBeginTransaction
	}

	_, err = stmt.ExecContext(ctx, form.Name, form.DrugID, form.Dose, form.AppliedAt, vaccinationId)

	if err != nil {
		repo.log.Info(err.Error())
		// log.Println("Code 2 ", errors.Is(err, my.ErrDupeKey))
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok {
			repo.log.Info("[INFO]", zap.Any("PGCode -> ", pgErr.Code))
			if pgErr.Code == "23505" {
				return errors.New("row already exists")
			} else {
				return ErrUpdatingRecord
			}
		}

	}
	if err = tx.Commit(); err != nil {
		return ErrCommitTransaction
	}
	return nil
}

func (repo repository) DeleteVaccinationItem(ctx context.Context, vaccinationId int) error {
	var query = `UPDATE vaccinations SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return ErrPrepapareQuery
	}
	defer func(stmt *sqlx.Stmt) {
		err := stmt.Close()
		if err != nil {
			repo.log.Error("failed to prepare statement", zap.Error(err))
		}
	}(stmt)
	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return ErrBeginTransaction
	}

	_, err = stmt.ExecContext(ctx, vaccinationId)

	if err != nil {
		repo.log.Info(err.Error())
		// log.Println("Code 2 ", errors.Is(err, my.ErrDupeKey))
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok {
			repo.log.Info("[INFO]", zap.Any("PG-CODE", pgErr.Code))
			if pgErr.Code == "23505" {
				return errors.New("row already exists")
			} else {
				return ErrDeletingRecord
			}
		}

	}
	if err = tx.Commit(); err != nil {
		return ErrCommitTransaction
	}
	return nil
}
