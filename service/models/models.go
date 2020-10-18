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
func (s *Session) FullDir(subs ...string) string {
	return path.Join(s.BaseDir, s.SubDir, path.Join(subs...))
}

// File represents a single file
type File struct {
	Name      string
	Ext       string
	DirPath   string
	CreatedAt time.Time
}

// NameWithExt returns the filename and extension of the associated file
func (f File) NameWithExt() string {
	return fmt.Sprintf("%s.%s", f.Name, f.Ext)
}

// FullPath returns the full path of the associated file
func (f File) FullPath() string {
	return path.Join(f.DirPath, f.NameWithExt())
}

// NewFile returns a new file object from the provided field values
func NewFile(name, ext, directory string, createdAt *time.Time) File {
	file := File{
		Name:    name,
		Ext:     ext,
		DirPath: directory,
	}
	if createdAt != nil {
		file.CreatedAt = *createdAt
	}
	return file
}

// Directory represents a single directory
type Directory struct {
	Name      string
	DirPath   string
	FileCount int
}

// FullPath returns the full path of the associated directory
func (d Directory) FullPath() string {
	return path.Join(d.DirPath, d.Name)
}
