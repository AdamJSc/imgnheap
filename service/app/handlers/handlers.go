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
		c.Templates().ExecuteTemplate(w, "index", views.IndexPage{Page: views.NewPage("Enter your directory", "", false)})
	}
}

func newImagesDirectoryHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessAgent := &domain.SessionAgent{SessionAgentInjector: c}

		// get directory path from request
		dirPath := r.FormValue("directory")
		if dirPath == "" {
			handleError(domain.BadRequestError{Err: errors.New("missing field: directory")}, c, w)
			return
		}

		// does directory exist?
		if !c.FileSystem().IsDirectory(dirPath) {
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
		redirect(w, "/catalog")
	}
}

func catalogMethodSelectionHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := c.FileSystem().GetFilesInDirectory(getDirPathFromRequest(r))
		if err != nil {
			handleError(err, c, w)
			return
		}

		imgFiles := domain.FilterFilesByExtensions(files, domain.ImgFileExts...)

		dirPath := getDirPathFromRequest(r)
		data := views.CatalogMethodSelectionPage{
			Page:            views.NewPage("Select your catalog method", dirPath, dirPath != ""),
			ImageFilesCount: len(imgFiles),
		}

		c.Templates().ExecuteTemplate(w, "catalog-method-selection", data)
	}
}

func resetHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// delete session cookie
		sessAgent := domain.SessionAgent{SessionAgentInjector: c}
		sessAgent.DeleteCookie(w)
		redirectToHome(w)
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

	data := views.ErrorPage{
		Page: views.NewPage(msg, "", true),
	}
	data.Error.Code = code
	data.Error.Message = msg
	data.Error.Detail = err.Error()

	c.Templates().ExecuteTemplate(w, "error", data)
}
