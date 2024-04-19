package models

import (
	"time"
)

type User struct {
	ID        int32     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	DeletedAt time.Time `json:"-" db:"deleted_at"`
}

// NewUser crea un nuevo usuario
func NewUser(name, password, email string) (*User, error) {
	user := &User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate valida al usuario
func (user *User) Validate() error {
	return nil
}
