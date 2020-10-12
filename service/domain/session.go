package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"imgnheap/service/app"
	"imgnheap/service/models"
	"net/http"
	"time"
)

const cookieName = "SESS_ID"

// SessionAgentInjector defines the injector behaviours for our SessionAgent
type SessionAgentInjector interface {
	app.FileSystemInjector
	app.KeyValStoreInjector
}

// SessionAgent represents our methods for interacting with sessions
type SessionAgent struct {
	SessionAgentInjector
}

// NewSessionFromRequestAndTimestamp generates a new session based on the provided request object and timestamp, and returns the session
func (s *SessionAgent) NewSessionFromRequestAndTimestamp(r *http.Request, ts time.Time) (*models.Session, error) {
	// get directory path from request
	dirPath := r.FormValue("directory")
	if dirPath == "" {
		return nil, BadRequestError{Err: errors.New("missing field: directory")}
	}

	// does directory exist?
	if !s.FileSystem().IsDirectory(dirPath) {
		return nil, ValidationError{Err: fmt.Errorf("not a directory: %s", dirPath)}
	}

	// generate new session token
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	sessToken := id.String()

	// create session object
	sess := &models.Session{
		Token:   sessToken,
		BaseDir: dirPath,
		SubDir:  fmt.Sprintf("processed%s", ts.Format("20060102150405")),
	}
	if err := s.KeyValStore().Write(sessToken, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

// GetSessionFromToken retrieves a Session object based on the provided token
func (s *SessionAgent) GetSessionFromToken(sessToken string) (*models.Session, error) {
	val, err := s.KeyValStore().Read(sessToken)
	if err != nil {
		return nil, err
	}

	sess, ok := val.(*models.Session)
	if !ok {
		return nil, fmt.Errorf("error token %s does not represents session object", sessToken)
	}

	return sess, nil
}

// WriteCookie writes the provided session as a cookie to the provided writer
func (s *SessionAgent) WriteCookie(sess *models.Session, w http.ResponseWriter) error {
	if sess == nil {
		return errors.New("session is nil")
	}

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: sess.Token,
		Path:  "/",
	})

	return nil
}

// DeleteCookie writes the removal of the provided session as a cookie to the provided writer
func (s *SessionAgent) DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Path:   "/",
		MaxAge: -1,
	})
}

// GetTokenFromCookie returns string value of session cookie, or empty string if missing
func (s *SessionAgent) GetTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie == nil {
		return ""
	}

	return cookie.Value
}
