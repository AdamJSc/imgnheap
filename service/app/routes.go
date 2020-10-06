package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	return r
}
