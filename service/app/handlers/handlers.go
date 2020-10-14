package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"imgnheap/service/app"
	"imgnheap/service/domain"
	"imgnheap/service/models"
	"imgnheap/service/views"
	"log"
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
		// define reusable method for writing output
		var writeResponse = func(data interface{}) {
			if err := c.Templates().ExecuteTemplate(w, "catalog-by-tag", data); err != nil {
				handleError(err, c, w)
			}
		}

		fsAgent := domain.FileSystemAgent{FileSystemAgentInjector: c}

		sess := getSessionFromRequest(r)
		if sess == nil {
			handleError(errors.New("session is nil"), c, w)
			return
		}
		dirPath := sess.BaseDir

		data := views.CatalogByTagPage{
			Page: views.NewPage("Catalog image by tag", dirPath, dirPath != ""),
		}

		// see if we have any more files that need to be processed
		files, err := fsAgent.GetFilesFromDirectoryByExtension(dirPath, domain.ImgFileExts...)
		if err != nil {
			handleError(err, c, w)
			return
		}
		data.ImageFilesCount = len(files)
		if data.ImageFilesCount == 0 {
			writeResponse(data)
			return
		}

		// get next file to be processed
		data.ImageFileName = files[0].FilenameWithExt()

		// TODO - get tags (subfolders) and file counts within each one

		writeResponse(data)
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
		file := models.NewFile(name, ext, sess.BaseDir, nil)

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

func renderFile(c app.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := getSessionFromRequest(r)
		if sess == nil {
			handleError(errors.New("session is nil"), c, w)
			return
		}

		var fileName string
		if err := routeParam(&fileName, "filename", r); err != nil {
			handleError(err, c, w)
			return
		}

		name, ext := domain.ParseNameAndExtensionFromFileName(fileName)
		file := models.NewFile(name, ext, sess.BaseDir, nil)

		fsAgent := domain.FileSystemAgent{FileSystemAgentInjector: c}

		err := fsAgent.Stream(file, w)
		if err != nil {
			w.WriteHeader(getResponseStatusFromError(err))
			log.Println(err)
		}
	}
}

// handleError handles the provided error and writes an appropriate error page
func handleError(err error, c app.Container, w http.ResponseWriter) {
	var msg string

	switch err.(type) {
	case domain.BadRequestError:
		msg = "Bad Request"
	case domain.NotFoundError:
		msg = "Not Found"
	case domain.ValidationError:
		msg = "Unprocessable Entity"
	default:
		msg = "Internal Server Error"
	}

	code := getResponseStatusFromError(err)

	data := views.ErrorPage{
		Page: views.NewPage(msg, "", true),
	}
	data.Error.Code = code
	data.Error.Message = msg
	data.Error.Detail = err.Error()

	w.WriteHeader(code)
	c.Templates().ExecuteTemplate(w, "error", data)
}

// getResponseStatusFromError returns a numeric response status code from the provided error
func getResponseStatusFromError(err error) int {
	switch err.(type) {
	case domain.BadRequestError:
		return http.StatusBadRequest
	case domain.NotFoundError:
		return http.StatusNotFound
	case domain.ValidationError:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// missingFieldError returns a new BadRequestError based on the provided field name
func missingFieldError(fieldName string) domain.BadRequestError {
	return domain.BadRequestError{
		Err: fmt.Errorf("missing field: %s", fieldName),
	}
}

// routeParam loads the value of the provided route parameter from the provided request object into the provided recipient variable
func routeParam(p *string, name string, r *http.Request) error {
	if p == nil {
		return errors.New("parameter is nil")
	}
	if r == nil {
		return errors.New("request is nil")
	}

	vars := mux.Vars(r)
	val, ok := vars[name]
	if !ok {
		return fmt.Errorf("param %s not found", name)
	}

	*p = val
	return nil
}
