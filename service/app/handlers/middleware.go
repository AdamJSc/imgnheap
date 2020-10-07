package handlers

import (
	"imgnheap/service/app"
	"imgnheap/service/domain"
	"net/http"
)

func sessionTokenIsValid(c app.Container) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// is session cookie set?
			sessAgent := domain.SessionAgent{SessionAgentInjector: c}
			if !sessAgent.IsSet(r) {
				// nope it isn't! redirect to home
				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusFound)
				return
			}

			// TODO - check that session cookie ID is valid
			// TODO - check that session cookie ID refers to a valid directory

			h.ServeHTTP(w, r)
		})
	}
}
