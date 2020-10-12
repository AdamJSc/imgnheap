package models

import "path"

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
