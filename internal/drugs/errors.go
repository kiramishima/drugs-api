package drugs

import "errors"

// Entity Errors
var (
	// Drugs
	InternalServerError   = errors.New("Error interno del servidor. Intente más tarde")
	ErrTimeout            = errors.New("context timeout")
	ErrPrepapareQuery     = errors.New("failed to prepare query")
	ErrDrugNotFound       = errors.New("No existe el medicamento")
	ErrServiceDrugs       = errors.New("service drug error")
	ErrExecuteStatement   = errors.New("failed to execute statement")
	ErrBeginTransaction   = errors.New("Falló al iniciar la transacción")
	ErrDuplicateDrug      = errors.New("Este medicamento ya existe")
	ErrCommitTransaction  = errors.New("failed to commit transaction")
	ErrRollback           = errors.New("failed to rollback")
	ErrInsertFailed       = errors.New("failed to insert new item")
	ErrNoRecords          = errors.New("No hay registros")
	ErrUpdatingRecord     = errors.New("failed to update record")
	ErrDeletingRecord     = errors.New("failed to delete record")
	ErrInvalidRequestBody = errors.New("El cuerpo de la petición es invalido")
)
