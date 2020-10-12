package handlers

import (
	"errors"
	"fmt"
	"imgnheap/service/app"
	"imgnheap/service/domain"
	"imgnheap/service/views"
	"net/http"
	"time"
)

func indexHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Templates().ExecuteTemplate(w, "index", views.IndexPage{Page: views.NewPage("Enter your directory", "", false)})
	}
}

func newSessionHandler(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessAgent := domain.SessionAgent{SessionAgentInjector: c}

		// save new session
		sess, err := sessAgent.NewSessionFromRequestAndTimestamp(r, time.Now())
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
			destDir := domain.GetDestinationDirFromFileTimestampAndSession(file, sess)
			if err := fsAgent.ProcessFileByCopy(file, sess.BaseDir, destDir); err != nil {
				handleError(err, c, w)
				return
			}
		}

		// TODO - execute template
		w.Write([]byte(fmt.Sprintf("processed %d files in %s", len(files), sess.FullDir())))
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
