package domain

import (
	"errors"
	"github.com/google/uuid"
	"imgnheap/service/app"
	"net/http"
)

const cookieName = "SESS_ID"

// Session defines a basic key/value session
type Session struct {
	Token   string
	DirPath string
}

// SessionAgentInjector defines the injector behaviours for our SessionAgent
type SessionAgentInjector interface {
	app.KeyValStoreInjector
}

// SessionAgent represents our methods for interacting with sessions
type SessionAgent struct {
	SessionAgentInjector
}

// NewSessionWithDirPath generates a new session, stores the provided directory path against it, and returns the token
func (s *SessionAgent) NewSessionWithDirPath(dirPath string) (*Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sessToken := id.String()

	if err := s.KeyValStore().Write(sessToken, dirPath); err != nil {
		return nil, err
	}

	return &Session{
		Token:   sessToken,
		DirPath: dirPath,
	}, nil
}

// GetSessionFromToken retrieves a Session object based on the provided token
func (s *SessionAgent) GetSessionFromToken(sessToken string) (*Session, error) {
	dirPath, err := s.KeyValStore().Read(sessToken)
	if err != nil {
		return nil, err
	}

	return &Session{
		Token:   sessToken,
		DirPath: dirPath,
	}, nil
}

// WriteCookie writes the provided session as a cookie to the provided writer
func (s *SessionAgent) WriteCookie(sess *Session, w http.ResponseWriter) error {
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
