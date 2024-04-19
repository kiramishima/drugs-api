package interfaces

import "net/http"

// DrugsHandlers interface
type DrugsHandlers interface {
	ListDrugsHandler(w http.ResponseWriter, req *http.Request)
	CreateDrugHandler(w http.ResponseWriter, req *http.Request)
	UpdateDrugHandler(w http.ResponseWriter, req *http.Request)
	DeleteDrugHandler(w http.ResponseWriter, req *http.Request)
}
