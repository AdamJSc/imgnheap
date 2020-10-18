package app

import (
	"html/template"
	"imgnheap/service/models"
)

// Container defines our app's container interface
type Container interface {
	TemplatesInjector
	KeyValStoreInjector
	FileSystemInjector
}

type TemplatesInjector interface{ Templates() *template.Template }
type KeyValStoreInjector interface{ KeyValStore() KeyValStore }
type FileSystemInjector interface{ FileSystem() FileSystem }

// KeyValStore defines operations for transacting with key/value storage
type KeyValStore interface {
	Read(key string) (interface{}, error)
	Write(key string, val interface{}) error
}

// FileSystem defines operations for transacting with a file system
type FileSystem interface {
	IsDirectory(path string) bool
	GetFilesInDirectory(path string) ([]models.File, error)
	GetDirectoriesInDirectory(path string) ([]models.Directory, error)
	GetContents(file models.File) ([]byte, error)
	Copy(file models.File, dest string) error
	Move(file models.File, dest string) error
}
