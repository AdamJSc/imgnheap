package domain

import (
	"imgnheap/service/app"
	"imgnheap/service/models"
	"os"
	"path/filepath"
	"strings"
)

// OsFileSystem defines the OS implementation of FileSystem
type OsFileSystem struct {
	app.FileSystem
}

// IsDirectory implements app.FileSystem.IsDirectory()
func (o *OsFileSystem) IsDirectory(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

// GetFilesInDirectory implements app.FileSystem.GetFilesInDirectory()
func (o *OsFileSystem) GetFilesInDirectory(dirPath string) ([]models.File, error) {
	var files []models.File

	if err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fname := info.Name()
		fnameParts := strings.Split(fname, ".")
		ext := fnameParts[len(fnameParts)-1]
		files = append(files, models.File{
			Name: fname,
			Ext:  ext,
		})

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

// filterFilesByExtensions returns the provided files filtered to retain only those whose extension matches one of the provided extensions
func filterFilesByExtensions(files []models.File, exts ...string) []models.File {
	var filtered []models.File

	for _, file := range files {
		if contains(exts, file.Ext) {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// contains returns true if the provided needle exists within the provided haystack, otherwise false
func contains(haystack []string, needle string) bool {
	for _, val := range haystack {
		if strings.ToLower(val) == strings.ToLower(needle) {
			return true
		}
	}

	return false
}
