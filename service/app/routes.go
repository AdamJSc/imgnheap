package app

import (
	"github.com/gorilla/mux"
)

func RegisterRouter(c Container) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler(c))

	return r
}
