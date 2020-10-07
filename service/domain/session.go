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
	Token string
	Value string
}

type SessionAgentInjector interface {
	app.KeyValStoreInjector
}

// SessionAgent represents our methods for interacting with sessions
type SessionAgent struct {
	SessionAgentInjector
}

// NewSessionWithValue generates a new session, stores the provided value against it, and returns the token
func (s *SessionAgent) NewSessionWithValue(val string) (*Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sessToken := id.String()

	if err := s.KeyValStore().Write(sessToken, val); err != nil {
		return nil, err
	}

	return &Session{
		Token: sessToken,
		Value: val,
	}, nil
}

// WriteCookie writes the provided session as a cookie to the provided writer
func (s *SessionAgent) WriteCookie(sess *Session, w http.ResponseWriter) error {
	if sess == nil {
		return errors.New("session is nil")
	}

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: sess.Value,
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

// IsSet returns true if session cookie is set, otherwise false
func (s *SessionAgent) IsSet(r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie == nil {
		return false
	}
	return true
}
