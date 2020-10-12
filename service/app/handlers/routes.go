package handlers

import (
	"github.com/gorilla/mux"
	"imgnheap/service/app"
	"net/http"
)

// RegisterRouter returns a new mux router with our handler routes attached
func RegisterRouter(c app.Container) *mux.Router {
	r := mux.NewRouter()

	// routes that require no session token
	r.HandleFunc("/", indexHandler(c)).Methods(http.MethodGet)
	r.HandleFunc("/", newSessionHandler(c)).Methods(http.MethodPost)

	// routes that require session token
	s := r.PathPrefix("").Subrouter()
	s.Use(addSessionToRequestContext(c))
	s.HandleFunc("/catalog", catalogMethodSelectionHandler(c)).Methods(http.MethodGet)
	s.HandleFunc("/catalog/by-date", processFilesByDateInFilename(c)).Methods(http.MethodPost)
	s.HandleFunc("/reset", resetHandler(c)).Methods(http.MethodPost)

	return r
}
