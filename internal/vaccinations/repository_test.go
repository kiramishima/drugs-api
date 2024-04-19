package vaccinations

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRepository_GetVaccinationsData(t *testing.T) {
	t.Parallel()
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

	repo := NewVaccinationRepository(sqlxDB, logger)

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

	var rows = sqlmock.NewRows([]string{"id", "name", "drug", "drug_id", "dose", "applied_at"}).FromCSVString("1,jhon wick,aspirina,1,5,2024-03-18 15:45:00")

	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnRows(rows)

		data, err := repo.GetVaccinationsData(ctx)
		t.Log(len(data), err)
		assert.NoError(t, err)
		assert.Equal(t, len(data), 2)
		assert.Equal(t, data[0].Name, "jhon wick")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No rows", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnError(sql.ErrNoRows)

		data, err := repo.GetVaccinationsData(ctx)
		t.Log(len(data), err)
		assert.NoError(t, err)
		assert.Equal(t, len(data), 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
