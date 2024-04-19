package vaccinations

import "errors"

// Entity Errors
var (
	// Drugs
	InternalServerError     = errors.New("Error interno del servidor. Intente más tarde")
	ErrTimeout              = errors.New("context timeout")
	ErrPrepapareQuery       = errors.New("Fallo al preparar el query")
	ErrVaccinationNotFound  = errors.New("Vacunación no encontrada")
	ErrServiceVaccination   = errors.New("Falla en el servicio vaccination")
	ErrExecuteStatement     = errors.New("Fallo al ejecutar la declaración SQL")
	ErrBeginTransaction     = errors.New("Fallo al iniciar la transacción")
	ErrDuplicateVaccination = errors.New("Registro existente")
	ErrCommitTransaction    = errors.New("Fallo al realizar el commit de la transacción")
	ErrRollback             = errors.New("Fallo al realizar el rollback")
	ErrInsertFailed         = errors.New("Fallo al insertar un nuevo registro")
	ErrNoRecords            = errors.New("No hay registros")
	ErrUpdatingRecord       = errors.New("Fallo al actualizar el registro")
	ErrDeletingRecord       = errors.New("Fallo al eliminar el registro")
	ErrInvalidRequestBody   = errors.New("El cuerpo de la petición es invalido")
)
