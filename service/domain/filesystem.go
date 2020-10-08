package domain

import (
	"imgnheap/service/app"
	"os"
)

// OsFileSystem defines the OS implementation of FileSystem
type OsFileSystem struct {
	app.FileSystem
}

// IsDir implements app.FileSystem.IsDir()
func (o *OsFileSystem) IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

// FileSystemAgentInjector defines the injector behaviours for our FileSystemAgent
type FileSystemAgentInjector interface {
	app.FileSystemInjector
}

// FileSystemAgent
type FileSystemAgent struct {
	FileSystemAgentInjector
}

func (f *FileSystemAgent) IsDir(path string) bool {
	return f.FileSystem().IsDir(path)
}
