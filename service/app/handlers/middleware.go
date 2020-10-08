package handlers

import (
	"context"
	"imgnheap/service/app"
	"imgnheap/service/domain"
	"net/http"
)

const ctxDirPathKey = "DIR_PATH"

// addDirPathToRequestContext provides a middleware method for adding the session's dir path to the request context
// otherwise, redirects current request to home if no valid session has been found
func addDirPathToRequestContext(c app.Container) func(http.Handler) http.Handler {
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

			if !c.FileSystem().IsDirectory(sess.DirPath) {
				// dir path stored by session token does not represent a valid directory
				sessAgent.DeleteCookie(w)
				redirectToHome(w)
				return
			}

			// add dir path to request context
			ctxWithDirPath := context.WithValue(r.Context(), ctxDirPathKey, sess.DirPath)
			h.ServeHTTP(w, r.WithContext(ctxWithDirPath))
		})
	}
}

// redirectToHome writes a redirection to the provided response writer
func redirectToHome(w http.ResponseWriter) {
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

// getDirPathFromRequest returns the dir path set on the request context by previous middleware
func getDirPathFromRequest(r *http.Request) string {
	val := r.Context().Value(ctxDirPathKey)
	if val == nil {
		return ""
	}

	dirPath, ok := val.(string)
	if !ok {
		return ""
	}

	return dirPath
}
