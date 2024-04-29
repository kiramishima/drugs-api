package auth

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

// implement auth repository
var _ interfaces.AuthRepository = (*repository)(nil)

// Repository struct
type repository struct {
	db  *sqlx.DB
	log *zap.Logger
}

// NewAuthRepository Creates a new instance of Repository
func NewAuthRepository(conn *sqlx.DB, logger *zap.Logger) *repository {
	return &repository{
		db:  conn,
		log: logger,
	}
}

func (repo repository) FindUserByCredentials(ctx context.Context, form *models.AuthForm) (*models.User, error) {
	var query = `SELECT id,
       	   name,
		   email,
		   password,
		   created_at,
		   updated_at
	FROM users
	WHERE email = $1`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, ErrPrepapareQuery
	}
	defer stmt.Close()

	u := &models.User{}

	row := stmt.QueryRowContext(ctx, form.Email)
	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	err = row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if createdAt.Valid {
		u.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		u.UpdatedAt = updatedAt.Time
	}

	return u, nil
}

func (repo repository) CreateAccount(ctx context.Context, form *models.RegisterForm) error {
	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		repo.log.Info(err.Error())
		return ErrBeginTransaction
	}
	defer tx.Rollback()

	// Prepare STMT
	var query = `INSERT INTO users(name, email, password) VALUES($1, $2, $3)`
	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return ErrPrepapareQuery
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, form.Name, form.Email, form.Password)

	if err != nil {
		repo.log.Info(err.Error())
		// log.Println("Code 2 ", errors.Is(err, my.ErrDupeKey))
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			repo.log.Info(pgErr.Code)
			if pgErr.Code == "23505" {
				return ErrUserExist
			} else {
				return ErrFailInsertUser
			}
		} else {
			return ErrFailInsertUser
		}

	}
	if err = tx.Commit(); err != nil {
		return ErrCommitTransaction
	}
	return nil
}
