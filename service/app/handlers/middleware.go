package handlers

import (
	"imgnheap/service/app"
	"imgnheap/service/domain"
	"net/http"
)

func sessionTokenIsValid(c app.Container) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessAgent := domain.SessionAgent{SessionAgentInjector: c}

			sessToken := sessAgent.GetTokenFromCookie(r)
			if sessToken == "" {
				// cookie is not even set! redirect to home
				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusFound)
				return
			}

			_, err := sessAgent.GetSessionFromToken(sessToken)
			if err != nil {
				// cookie value does not represent a valid session token! redirect to home
				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusFound)
				return
			}

			// TODO - check that session cookie ID refers to a valid directory

			h.ServeHTTP(w, r)
		})
	}
}
