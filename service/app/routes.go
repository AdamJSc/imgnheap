package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterRouter returns a new mux router with our handler routes attached
func RegisterRouter(c Container) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler(c)).Methods(http.MethodGet)
	r.HandleFunc("/", newImagesDirectoryHandler(c)).Methods(http.MethodPost)

	return r
}
