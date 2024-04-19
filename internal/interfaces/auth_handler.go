package interfaces

import "net/http"

type AuthHandlers interface {
	SignUpHandler(w http.ResponseWriter, req *http.Request)
	LoginHandler(w http.ResponseWriter, req *http.Request)
}
