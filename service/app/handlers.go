package app

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"imgnheap/service/domain"
	"imgnheap/service/views"
	"net/http"
	"os"
)

const sessionCookieName = "SESS_ID"

func indexHandler(c Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Templates().ExecuteTemplate(w, "index", views.IndexPage{Page: views.Page{Title: "Select Directory"}})
	}
}

func newImagesDirectoryHandler(c Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get directory path
		dirPath := r.FormValue("directory")
		if dirPath == "" {
			handleError(domain.BadRequestError{Err: errors.New("missing field: directory")}, c, w)
			return
		}

		// does directory exist?
		fi, err := os.Stat(dirPath)
		if err != nil {
			handleError(domain.ValidationError{Err: fmt.Errorf("directory does not exist: %s", dirPath)}, c, w)
			return
		}
		if !fi.IsDir() {
			handleError(domain.ValidationError{Err: fmt.Errorf("not a directory: %s", dirPath)}, c, w)
			return
		}

		// generate a new session ID
		sessID, err := uuid.NewRandom()
		if err != nil {
			handleError(err, c, w)
			return
		}

		// TODO - save path against session ID

		// write session cookie
		http.SetCookie(w, &http.Cookie{
			Name:  sessionCookieName,
			Value: sessID.String(),
			Path:  "/",
		})

		// write location header
		// TODO - amend location
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
	}
}

// handleError handles the provided error and writes an appropriate error page
func handleError(err error, c Container, w http.ResponseWriter) {
	var code int
	var msg string

	switch err.(type) {
	case domain.BadRequestError:
		code = http.StatusBadRequest
		msg = "Bad Request"
	case domain.ValidationError:
		code = http.StatusUnprocessableEntity
		msg = "Unprocessable Entity"
	default:
		code = http.StatusInternalServerError
		msg = "Internal Server Error"
	}

	var data views.ErrorPage
	data.Title = msg
	data.Error.Code = code
	data.Error.Message = msg
	data.Error.Detail = err.Error()

	c.Templates().ExecuteTemplate(w, "error", data)
}
