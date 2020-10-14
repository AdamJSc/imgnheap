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

// NewFile returns a new file object from the provided field values
func NewFile(name, ext, directory string, createdAt *time.Time) File {
	file := File{
		Name:      name,
		Ext:       ext,
		Directory: directory,
	}
	if createdAt != nil {
		file.CreatedAt = *createdAt
	}
	return file
}
