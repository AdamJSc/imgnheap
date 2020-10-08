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
			fsAgent := domain.FileSystemAgent{FileSystemAgentInjector: c}

			sessToken := sessAgent.GetTokenFromCookie(r)
			if sessToken == "" {
				// cookie is not set
				redirectToHome(w)
				return
			}

			sess, err := sessAgent.GetSessionFromToken(sessToken)
			if err != nil {
				// cookie value does not represent a valid session token
				redirectToHome(w)
				return
			}

			if !fsAgent.IsDir(sess.DirPath) {
				// dir path stored by session token does not represent a valid directory
				redirectToHome(w)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

// redirectToHome writes a redirection to the provided response writer
func redirectToHome(w http.ResponseWriter) {
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}
