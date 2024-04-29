package drugs

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/models"
	"testing"
	"time"
)

func TestRepository_GetDrugsData(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("close db", zap.Error(err))
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()

	repo := NewDrugRepository(sqlxDB, logger)

	var query = `SELECT id, name, approved, min_dose, max_dose, available_at FROM drugs WHERE deleted_at IS NULL`

	var rows = sqlmock.NewRows([]string{"id", "name", "approved", "min_dose", "max_dose", "available_at"}).FromCSVString("1,aspirina,true,1,5,2024-05-05 00:00:00\n2,cafiaspirina,true,2,5,2024-05-05 00:00:00")

	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnRows(rows)

		data, err := repo.GetDrugsData(ctx)
		t.Log(len(data), err)
		assert.NoError(t, err)
		assert.Equal(t, len(data), 2)
		assert.Equal(t, data[0].Name, "aspirina")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No rows", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnError(sql.ErrNoRows)

		data, err := repo.GetDrugsData(ctx)
		t.Log(len(data), err)
		assert.NoError(t, err)
		assert.Equal(t, len(data), 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_GetDrugItemByID(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("", zap.Error(err))
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()

	repo := NewDrugRepository(sqlxDB, logger)

	var query = `SELECT id, name, approved, min_dose, max_dose, available_at FROM drugs 
    WHERE deleted_at IS NULL AND id = $1`

	var rows = sqlmock.NewRows([]string{"id", "name", "approved", "min_dose", "max_dose", "available_at"}).FromCSVString("1,aspirina,true,1,5,2024-05-05 00:00:00")

	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnRows(rows)

		data, err := repo.GetDrugItemByID(ctx, 1)
		t.Log(data, err)
		assert.NoError(t, err)
		assert.Equal(t, data.Name, "aspirina")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No row", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetDrugItemByID(ctx, 1)
		// t.Log(len(data), err)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_CreateNewDrugItem(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("", zap.Error(err))
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()

	repo := NewDrugRepository(sqlxDB, logger)

	var query = `INSERT INTO drugs (name, approved, min_dose, max_dose, available_at)
	VALUES ($1, $2, $3, $4, CAST($5 AS TIMESTAMP))`

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

	t.Run("Insert is OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(10)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(item.Name, item.Approved, item.MinDose, item.MaxDose, item.AvailableAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.CreateNewDrugItem(ctx, item)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Duplicate record", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WillReturnError(&pgconn.PgError{
				Code: "23505", // Duplicate key error code
			})

		mock.ExpectRollback()

		err := repo.CreateNewDrugItem(ctx, item)
		t.Log(err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrDuplicateDrug.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail begin transaction", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query)

		mock.ExpectBegin().WillReturnError(ErrBeginTransaction)

		mock.ExpectExec(query).
			WithArgs(item.Name, item.Approved, item.MinDose, item.MaxDose, item.AvailableAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectRollback()

		err := repo.CreateNewDrugItem(ctx, item)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrBeginTransaction.Error())
		assert.Error(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_UpdateDrugItem(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("", zap.Error(err))
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()

	repo := NewDrugRepository(sqlxDB, logger)

	var query = `UPDATE drugs SET name = $1, approved = $2, min_dose = $3, max_dose = $4, available_at = $5 WHERE id = $6`

	var name = "Aspirina"
	var approved = true
	var minDose = 1
	var maxDose = 2
	var availableAt = time.Now()
	var item = &models.Drug{
		ID:          1,
		Name:        name,
		Approved:    approved,
		MinDose:     minDose,
		MaxDose:     maxDose,
		AvailableAt: availableAt,
	}

	t.Run("Updated is OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(item.Name, item.Approved, item.MinDose, item.MaxDose, item.AvailableAt, item.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.UpdateDrugItem(ctx, 1, item)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Updated Not Existing Item", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(item.Name, item.Approved, item.MinDose, item.MaxDose, item.AvailableAt, item.ID).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		err := repo.UpdateDrugItem(ctx, 1, item)
		t.Log(err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrDrugNotFound.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_DeleteDrugItem(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("", zap.Error(err))
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()

	repo := NewDrugRepository(sqlxDB, logger)

	var query = `UPDATE drugs SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	t.Run("Deleted is OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.DeleteDrugItem(ctx, 1)
		t.Log(err)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail Deleting Not Existing Item", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		err := repo.DeleteDrugItem(ctx, 1)
		t.Log(err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrDrugNotFound.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
