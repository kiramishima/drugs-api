package drugs

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/interfaces"
	"kiramishima/ionix/internal/models"
)

// implement drug repository
var _ interfaces.DrugRepository = (*repository)(nil)

// Repository struct
type repository struct {
	db  *sqlx.DB
	log *zap.SugaredLogger
}

// NewDrugRepository Creates a new instance of Repository
func NewDrugRepository(conn *sqlx.DB, logger *zap.SugaredLogger) *repository {
	return &repository{
		db:  conn,
		log: logger,
	}
}

// GetDrugsData gets data from drugs table
func (repo repository) GetDrugsData(ctx context.Context) ([]*models.Drug, error) {
	var query = `SELECT id, name, approved, min_dose, max_dose, available_at FROM drugs WHERE deleted_at IS NULL`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, ErrPrepapareQuery
	}
	defer stmt.Close()

	var list = make([]*models.Drug, 0)

	rows, err := stmt.QueryxContext(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return list, nil
		} else {
			return list, ErrExecuteStatement
		}
	}
	for rows.Next() {
		var availableAt sql.NullTime
		var item = &models.Drug{}
		err = rows.Scan(&item.ID, &item.Name, &item.Approved, &item.MinDose, &item.MaxDose, &availableAt)
		repo.log.Info(item)
		if errors.Is(err, sql.ErrNoRows) {
			break
		}

		if availableAt.Valid {
			item.AvailableAt = availableAt.Time
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo repository) GetDrugItemByID(ctx context.Context, drugId int) (*models.Drug, error) {
	var query = `SELECT id, name, approved, min_dose, max_dose, available_at FROM drugs 
    WHERE deleted_at IS NULL AND id = $1`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, ErrPrepapareQuery
	}
	defer stmt.Close()

	rows := stmt.QueryRowContext(ctx, drugId)

	var availableAt sql.NullTime
	var item = &models.Drug{}
	err = rows.Scan(&item.ID, &item.Name, &item.Approved, &item.MinDose, &item.MaxDose, &availableAt)
	repo.log.Info(item)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDrugNotFound
	}

	if availableAt.Valid {
		item.AvailableAt = availableAt.Time
	}

	return item, nil
}

func (repo repository) CreateNewDrugItem(ctx context.Context, form *models.DrugForm) error {
	repo.log.Info(form)
	var query = `INSERT INTO drugs (name, approved, min_dose, max_dose, available_at)
	VALUES ($1, $2, $3, $4, CAST($5 AS TIMESTAMP))`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return ErrPrepapareQuery
	}
	defer stmt.Close()

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return ErrBeginTransaction
	}

	_, err = stmt.ExecContext(ctx, form.Name, form.Approved, form.MinDose, form.MaxDose, form.AvailableAt)

	if err != nil {
		repo.log.Info(err.Error())
		tx.Rollback()
		// log.Println("Code 2 ", errors.Is(err, my.ErrDupeKey))
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			repo.log.Info(pgErr.Code)
			if pgErr.Code == "23505" {
				return ErrDuplicateDrug
			} else {
				return ErrInsertFailed
			}
		}

	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return ErrCommitTransaction
	}
	return nil
}

func (repo repository) UpdateDrugItem(ctx context.Context, drugId int, form *models.Drug) error {
	repo.log.Info(form)
	var query = `UPDATE drugs SET name = $1, approved = $2, min_dose = $3, max_dose = $4, available_at = $5 WHERE id = $6`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return ErrPrepapareQuery
	}
	defer stmt.Close()

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return ErrBeginTransaction
	}

	_, err = stmt.ExecContext(ctx, form.Name, form.Approved, form.MinDose, form.MaxDose, form.AvailableAt, drugId)

	if err != nil {
		tx.Rollback()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrDrugNotFound
		default:
			return ErrUpdatingRecord
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return ErrCommitTransaction
	}
	return nil
}

func (repo repository) DeleteDrugItem(ctx context.Context, drugId int) error {
	var query = `UPDATE drugs SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return ErrPrepapareQuery
	}
	defer stmt.Close()

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return ErrBeginTransaction
	}

	_, err = stmt.ExecContext(ctx, drugId)

	if err != nil {
		tx.Rollback()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrDrugNotFound
		default:
			return ErrUpdatingRecord
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return ErrCommitTransaction
	}
	return nil
}
