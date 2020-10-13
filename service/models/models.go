package models

import (
	"fmt"
	"path"
	"time"
)

// Session defines a basic session
type Session struct {
	Token   string
	BaseDir string
	SubDir  string
}

// FullDir returns the full directory stored by the Session
func (s *Session) FullDir() string {
	return path.Join(s.BaseDir, s.SubDir)
}

// File represents a single file
type File struct {
	Name      string
	Ext       string
	Directory string
	CreatedAt time.Time
}

// FilenameWithExt returns the filename and extension of the associated file
func (f File) FilenameWithExt() string {
	return fmt.Sprintf("%s.%s", f.Name, f.Ext)
}

// FullPath returns the full path of the associated file
func (f File) FullPath() string {
	return path.Join(f.Directory, f.FilenameWithExt())
}
