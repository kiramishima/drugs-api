package auth

import "errors"

// Entity Errors
var (
	// User
	ErrPrepapareQuery    = errors.New("Falló al preparar la consulta")
	ErrFieldWithSpaces   = errors.New("username and password can't have spaces")
	ErrShortPassword     = errors.New("password shorter than 6 characters")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrMissingPassword   = errors.New("El campo password es requerido")
	ErrMissingEmail      = errors.New("El campo email es requerido")
	ErrInvalidEmail      = errors.New("El campo email es invalido")
	ErrUserNotFound      = errors.New("Usuario no existe")
	ErrUserExist         = errors.New("Ya existe una cuenta con esta email")
	ErrFailInsertUser    = errors.New("Falló al registrar nuevo usuario")
	ErrServiceAuth       = errors.New("Falló el servicio auth")
	ErrBeginTransaction  = errors.New("Error al iniciar la transacción")
	ErrCommitTransaction = errors.New("Error al realizar el commit")
)
