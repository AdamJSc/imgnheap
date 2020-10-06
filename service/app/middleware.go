package app

import "net/http"

func sessionTokenIsValid(c Container) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// is session cookie set?
			_, err := r.Cookie(sessionCookieName)
			if err != nil {
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
