package handlers

import (
	"context"
	"imgnheap/service/app"
	"imgnheap/service/domain"
	"imgnheap/service/models"
	"net/http"
)

const ctxSessionKey = "CTX_SESSION"

// addSessionToRequestContext provides a middleware method for adding the session to the request context
// otherwise, redirects current request to home if no valid session has been found
func addSessionToRequestContext(c app.Container) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessAgent := domain.SessionAgent{SessionAgentInjector: c}

			sessToken := sessAgent.GetTokenFromCookie(r)
			if sessToken == "" {
				// cookie is not set
				redirectToHome(w)
				return
			}

			sess, err := sessAgent.GetSessionFromToken(sessToken)
			if err != nil {
				// cookie value does not represent a valid session token
				sessAgent.DeleteCookie(w)
				redirectToHome(w)
				return
			}

			if !c.FileSystem().IsDirectory(sess.BaseDir) {
				// dir path stored by session token does not represent a valid directory
				sessAgent.DeleteCookie(w)
				redirectToHome(w)
				return
			}

			// add session to request context
			ctxWithDirPath := context.WithValue(r.Context(), ctxSessionKey, sess)
			h.ServeHTTP(w, r.WithContext(ctxWithDirPath))
		})
	}
}

// redirectToHome writes a redirection to the provided response writer
func redirectToHome(w http.ResponseWriter) {
	redirect(w, "/")
}

// redirect writes to the provided response writer a redirection to the provided location
func redirect(w http.ResponseWriter, location string) {
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
}

// getSessionFromRequest returns the session set on the request context by previous middleware
func getSessionFromRequest(r *http.Request) *models.Session {
	val := r.Context().Value(ctxSessionKey)
	if val == nil {
		return nil
	}

	sess, ok := val.(*models.Session)
	if !ok {
		return nil
	}

	return sess
}
