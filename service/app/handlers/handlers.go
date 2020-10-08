package handlers

import (
	"errors"
	"fmt"
	"imgnheap/service/app"
	"imgnheap/service/domain"
	"imgnheap/service/views"
	"net/http"
)

func indexHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Templates().ExecuteTemplate(w, "index", views.IndexPage{Page: views.Page{Title: "Enter your directory"}})
	}
}

func newImagesDirectoryHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessAgent := &domain.SessionAgent{SessionAgentInjector: c}
		fsAgent := domain.FileSystemAgent{FileSystemAgentInjector: c}

		// get directory path from request
		dirPath := r.FormValue("directory")
		if dirPath == "" {
			handleError(domain.BadRequestError{Err: errors.New("missing field: directory")}, c, w)
			return
		}

		// does directory exist?
		if !fsAgent.IsDir(dirPath) {
			handleError(domain.ValidationError{Err: fmt.Errorf("not a directory: %s", dirPath)}, c, w)
			return
		}

		// save new session
		sess, err := sessAgent.NewSessionWithDirPath(dirPath)
		if err != nil {
			handleError(err, c, w)
			return
		}

		// write cookie
		if err := sessAgent.WriteCookie(sess, w); err != nil {
			handleError(err, c, w)
			return
		}

		// redirect to next step
		w.Header().Set("Location", "/catalog")
		w.WriteHeader(http.StatusFound)
	}
}

func catalogMethodSelectionHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO - get image files count from directory
		// TODO - ensure there's at least one file to process

		imageFilesCount := 0

		var data views.CatalogMethodSelectionPage
		data.Title = "Select your catalog method"
		data.ImageFilesCount = imageFilesCount

		c.Templates().ExecuteTemplate(w, "catalog-method-selection", data)
	}
}

func resetHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// delete session cookie
		sessAgent := domain.SessionAgent{SessionAgentInjector: c}
		sessAgent.DeleteCookie(w)

		// redirect to home
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
	}
}

// handleError handles the provided error and writes an appropriate error page
func handleError(err error, c app.Container, w http.ResponseWriter) {
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
