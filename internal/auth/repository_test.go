package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"kiramishima/ionix/internal/models"
	"testing"
	"time"
)

func TestRepository_CreateAccount(t *testing.T) {
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

	repo := NewAuthRepository(sqlxDB, logger)

	var query = `INSERT INTO users(name, email, password) VALUES($1, $2, $3)`

	var item = &models.RegisterForm{
		Name:     "Jhon Wick",
		Email:    "jhonwick@gmail.com",
		Password: "123456",
	}
	item.Password = item.Hash256Password(item.Password)

	t.Run("Insert is OK", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(c, time.Duration(10)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(item.Name, item.Email, item.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.CreateAccount(ctx, item)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail start transactionm", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin().WillReturnError(ErrBeginTransaction)

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(item.Name, item.Email, item.Password).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err := repo.CreateAccount(ctx, item)
		t.Log("err", err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrBeginTransaction.Error())
		assert.Error(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail prepare query", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			WillReturnError(ErrPrepapareQuery)

		mock.ExpectRollback()

		err := repo.CreateAccount(ctx, item)
		t.Log("err", err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrPrepapareQuery.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Printf("unmet expectation error: %s", err)
		}
	})

	t.Run("Duplicate user", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(item.Name, item.Email, item.Password).
			WillReturnError(&pgconn.PgError{
				Code: "23505", // Duplicate key error code
			})

		mock.ExpectRollback()

		err := repo.CreateAccount(ctx, item)
		t.Log("err", err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrUserExist.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail insert user", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(item.Name, item.Email, item.Password).
			WillReturnError(ErrFailInsertUser)

		mock.ExpectRollback()

		err := repo.CreateAccount(ctx, item)
		t.Log("err", err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrFailInsertUser.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail commit transactionm", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectBegin()

		mock.ExpectPrepare(query).
			ExpectExec().
			WithArgs(item.Name, item.Email, item.Password).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit().WillReturnError(ErrCommitTransaction)

		err := repo.CreateAccount(ctx, item)
		t.Log("err", err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrCommitTransaction.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_FindUserByCredentials(t *testing.T) {
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

	repo := NewAuthRepository(sqlxDB, logger)

	var query = `SELECT id,
       	   name,
		   email,
		   password,
		   created_at,
		   updated_at
	FROM users
	WHERE email = $1`

	var item = &models.User{
		ID:        1,
		Name:      "Jhon Wick",
		Email:     "jhonwick@gmail.com",
		Password:  "123456",
		CreatedAt: time.Now(),
	}

	var form = &models.AuthForm{
		Email:    "jhonwick@gmail.com",
		Password: "123456",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
		AddRow(item.ID, item.Name, item.Email, item.Password, item.CreatedAt, nil)

	t.Run("OK", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			ExpectQuery().
			WithArgs(form.Email).
			WillReturnRows(rows)

		user, err := repo.FindUserByCredentials(ctx, form)
		t.Log(user, "err ", err)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, item.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail prepare query", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			WillReturnError(ErrPrepapareQuery)

		_, err := repo.FindUserByCredentials(ctx, form)
		t.Log("err", err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrPrepapareQuery.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not exist user", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
		defer cancel()

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnError(sql.ErrNoRows)

		_, err := repo.FindUserByCredentials(ctx, form)
		t.Log("err", err)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrUserNotFound.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
