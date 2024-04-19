package vaccinations

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/models"
	"testing"
	"time"
)

func TestRepository_RegisterAssign(t *testing.T) {
	zlog, _ := zap.NewProduction()
	logger := zlog.Sugar()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err)
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()
	ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
	defer cancel()

	repo := NewCreditRepository(sqlxDB, logger)

	var query = `INSERT INTO credit_assigns(invest, credit_300, credit_500, credit_700, status) 
		VALUES($1, $2, $3, $4, $5)`

	var item = &models.Credit{
		Invest:    3000,
		Credit300: 2,
		Credit500: 2,
		Credit700: 2,
		Status:    1,
	}

	t.Run("Insert is OK", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(query).
			WithArgs(item.Invest, item.Credit300, item.Credit500, item.Credit700, item.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.RegisterAssign(ctx, item)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail start transactionm", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("Error al iniciar la transacción"))

		mock.ExpectExec(query).
			WithArgs(item.Invest, item.Credit300, item.Credit500, item.Credit700, item.Status).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err := repo.RegisterAssign(ctx, item)
		t.Log("err", err)
		assert.Error(t, err)
		assert.EqualError(t, err, "Error al iniciar la transacción")
		assert.Error(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_SumUp(t *testing.T) {
	zlog, _ := zap.NewProduction()
	logger := zlog.Sugar()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err)
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()
	ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
	defer cancel()

	repo := NewCreditRepository(sqlxDB, logger)

	var query = `SELECT total,
		total_sucess, 
		total_fails, 
		avg_total_success_inv,
		avg_total_fail_inv 
	FROM statistics`

	var item = &models.Stats{
		TotalAssigns:        26,
		TotalSuccessAssigns: 16,
		TotalFailAssigns:    10,
		AVGSuccessAssigns:   70.13,
		AVGFailAssigns:      29.87,
	}

	rows := sqlmock.NewRows([]string{"total", "total_sucess", "total_fails", "avg_total_success_inv", "avg_total_fail_inv"}).
		AddRow(&item.TotalAssigns, &item.TotalSuccessAssigns, &item.TotalFailAssigns, &item.AVGSuccessAssigns, &item.AVGFailAssigns)

	t.Run("Ok", func(t *testing.T) {

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnRows(rows)

		stats, err := repo.SumUp(ctx)
		assert.NoError(t, err)
		assert.Equal(t, stats.TotalAssigns, item.TotalAssigns)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No records", func(t *testing.T) {
		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnError(sql.ErrNoRows)

		stats, err := repo.SumUp(ctx)
		assert.Error(t, err)
		// t.Log(stats, err)
		assert.Equal(t, stats.TotalAssigns, int64(0))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
