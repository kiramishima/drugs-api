package interfaces

import "net/http"

// VaccinationsHandlers interface
type VaccinationsHandlers interface {
	ListVaccinationsHandler(w http.ResponseWriter, req *http.Request)
	CreateVaccinationHandler(w http.ResponseWriter, req *http.Request)
	UpdateVaccinationHandler(w http.ResponseWriter, req *http.Request)
	DeleteVaccinationHandler(w http.ResponseWriter, req *http.Request)
}
