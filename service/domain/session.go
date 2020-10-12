package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"imgnheap/service/app"
	"net/http"
	"path"
	"time"
)

const cookieName = "SESS_ID"

// Session defines a basic key/value session
type Session struct {
	Token   string
	BaseDir string
	SubDir  string
}

// FullDir returns the full directory stored by the Session
func (s *Session) FullDir() string {
	return path.Join(s.BaseDir, s.SubDir)
}

// SessionAgentInjector defines the injector behaviours for our SessionAgent
type SessionAgentInjector interface {
	app.KeyValStoreInjector
}

// SessionAgent represents our methods for interacting with sessions
type SessionAgent struct {
	SessionAgentInjector
}

// NewSessionWithDirPathAndTimestamp generates a new session, stores the provided directory path and timestamp against it, and returns the token
func (s *SessionAgent) NewSessionWithDirPathAndTimestamp(dirPath string, ts time.Time) (*Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sessToken := id.String()

	sess := &Session{
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
func (s *SessionAgent) GetSessionFromToken(sessToken string) (*Session, error) {
	val, err := s.KeyValStore().Read(sessToken)
	if err != nil {
		return nil, err
	}

	sess, ok := val.(*Session)
	if !ok {
		return nil, fmt.Errorf("error token %s does not represents session object", sessToken)
	}

	return sess, nil
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
