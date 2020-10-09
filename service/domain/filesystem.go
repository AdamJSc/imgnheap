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

		fileName, ext := ParseFileNameAndExtensionFromInfo(info)

		files = append(files, models.File{
			Name:      fileName,
			Ext:       ext,
			CreatedAt: info.ModTime(),
		})

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

// ParseFileNameAndExtensionFromInfo returns the filename and extension from the provided file info object
func ParseFileNameAndExtensionFromInfo(info os.FileInfo) (string, string) {
	var ext string
	fileName := info.Name()
	fileNameParts := strings.Split(fileName, ".")

	if len(fileNameParts) > 1 {
		fileName = strings.Join(fileNameParts[0:len(fileNameParts)-1], ".")
		ext = fileNameParts[len(fileNameParts)-1]
	}

	return fileName, ext
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
