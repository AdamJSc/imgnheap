package handlers

import (
	"errors"
	"fmt"
	"imgnheap/service/app"
	"imgnheap/service/domain"
	"imgnheap/service/models"
	"imgnheap/service/views"
	"net/http"
	"time"
)

func indexHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := views.IndexPage{Page: views.NewPage("Enter your directory", "", false)}

		if err := c.Templates().ExecuteTemplate(w, "index", data); err != nil {
			handleError(err, c, w)
		}
	}
}

func newSessionHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get directory path from request
		dirPath := r.FormValue("directory")
		if dirPath == "" {
			handleError(missingFieldError("directory"), c, w)
			return
		}

		sessAgent := domain.SessionAgent{SessionAgentInjector: c}

		// save new session
		sess, err := sessAgent.NewSessionFromDirectoryAndTimestamp(dirPath, time.Now())
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
		fsAgent := domain.FileSystemAgent{FileSystemAgentInjector: c}

		sess := getSessionFromRequest(r)
		if sess == nil {
			handleError(errors.New("session is nil"), c, w)
			return
		}
		dirPath := sess.BaseDir

		imgFiles, err := fsAgent.GetFilesFromDirectoryByExtension(dirPath, domain.ImgFileExts...)
		if err != nil {
			handleError(err, c, w)
			return
		}

		data := views.CatalogMethodSelectionPage{
			Page:            views.NewPage("Select your catalog method", dirPath, dirPath != ""),
			ImageFilesCount: len(imgFiles),
		}

		if err := c.Templates().ExecuteTemplate(w, "catalog-method-selection", data); err != nil {
			handleError(err, c, w)
		}
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

func processFilesByDateInFilename(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := getSessionFromRequest(r)
		if sess == nil {
			handleError(errors.New("session is nil"), c, w)
			return
		}

		fsAgent := domain.FileSystemAgent{FileSystemAgentInjector: c}

		files, err := fsAgent.GetFilesFromDirectoryByExtension(sess.BaseDir, domain.ImgFileExts...)
		if err != nil {
			handleError(err, c, w)
			return
		}

		for _, file := range files {
			destDir := domain.GetDestinationDirByDate(file, sess)
			if err := fsAgent.ProcessFileByCopy(file, destDir); err != nil {
				handleError(err, c, w)
				return
			}
		}

		// TODO - execute template
		w.Write([]byte(fmt.Sprintf("processed %d files in %s", len(files), sess.FullDir())))
	}
}

func catalogByTag(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO - implement handler
		w.Write([]byte("thing coming soon..."))
	}
}

func processFileByTag(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := getSessionFromRequest(r)
		if sess == nil {
			handleError(errors.New("session is nil"), c, w)
			return
		}

		// get filename from request
		fileName := r.FormValue("file_name")
		if fileName == "" {
			handleError(missingFieldError("file_name"), c, w)
			return
		}

		// get tag from request
		tag := r.FormValue("tag")
		if tag == "" {
			handleError(missingFieldError("tag"), c, w)
			return
		}

		// instantiate file object
		name, ext := domain.ParseNameAndExtensionFromFileName(fileName)
		file := models.NewFile(name, ext, sess.BaseDir, time.Time{})

		fsAgent := domain.FileSystemAgent{FileSystemAgentInjector: c}

		// do the copy bit...
		destDir := domain.GetDestinationDirByTag(sess, tag)
		if err := fsAgent.ProcessFileByCopy(file, destDir); err != nil {
			handleError(err, c, w)
			return
		}

		// redirect to control panel
		redirect(w, "/catalog/by-tag")
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

// missingFieldError returns a new BadRequestError based on the provided field name
func missingFieldError(fieldName string) domain.BadRequestError {
	return domain.BadRequestError{
		Err: fmt.Errorf("missing field: %s", fieldName),
	}
}
