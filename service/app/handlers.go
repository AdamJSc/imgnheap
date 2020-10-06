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
		// is session cookie set?
		_, err := r.Cookie(sessionCookieName)
		if err != nil {
			// nope it isn't!
			c.Templates().ExecuteTemplate(w, "index", views.IndexPage{Page: views.Page{Title: "Enter your directory"}})
			return
		}

		// TODO - check that session cookie ID is valid
		// TODO - check that session cookie ID refers to a valid directory
		// TODO - get image files count from directory
		// TODO - ensure there's at least one file to process

		imageFilesCount := 0

		var data views.CatalogMethodSelectionPage
		data.Title = "Select your catalog method"
		data.ImageFilesCount = imageFilesCount

		c.Templates().ExecuteTemplate(w, "catalog-method-selection", data)
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

		// redirect to home
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
	}
}

func resetHandler(c Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// delete session cookie
		http.SetCookie(w, &http.Cookie{
			Name:       sessionCookieName,
			Path:       "/",
			MaxAge:     -1,
		})

		// redirect to home
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
