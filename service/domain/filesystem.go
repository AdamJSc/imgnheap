package domain

import (
	"imgnheap/service/app"
	"imgnheap/service/models"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ImgFileExts = []string{
	"png",
	"jpg",
	"jpeg",
	"mp4",
}

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

// FileSystemAgentInjector defines the injector behaviours for our FileSystemAgent
type FileSystemAgentInjector interface {
	app.FileSystemInjector
}

// FileSystemAgent encapsulates all of our filesystem-related operations
type FileSystemAgent struct {
	FileSystemAgentInjector
}

// GetFilesFromDirectoryByExtension returns a slice of the files present within the provided directory path
func (i *FileSystemAgent) GetFilesFromDirectoryByExtension(dir string, exts ...string) ([]models.File, error) {
	files, err := i.FileSystem().GetFilesInDirectory(dir)
	if err != nil {
		return nil, err
	}

	var filtered []models.File

	for _, file := range files {
		if contains(exts, file.Ext) {
			filtered = append(filtered, file)
		}
	}

	return filtered, nil
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

// ParseTimestampFromFile attempts to parse a timestamp from the provided file
func ParseTimestampFromFile(file models.File) time.Time {
	// define potential date-based file naming patterns
	tsLayouts := []string{
		"20060102150405",
		"20060102_150405",
		"20060102-150405",
		"Screenshot_20060102150405",
		"Screenshot_20060102_150405",
		"Screenshot_20060102-150405",
		"Screenshot 2006-01-02 at 15.04.05",
	}

	var tsFromFileNameVariations = func(fileName string) []string {
		variations := []string{
			// original filename
			fileName,
		}

		// var1 = "Screenshot_<timestamp>_<some_app_name>"
		if strings.HasPrefix(fileName, "Screenshot_") {
			// suffix = "<timestamp>_<some_app_name>"
			suffix := strings.SplitN(fileName, "Screenshot_", 2)[1]
			suffixParts := strings.Split(suffix, "_")

			var1 := strings.Join(suffixParts[0:len(suffixParts)-1], "_")
			variations = append(variations, var1)
		}

		return variations
	}

	for _, layout := range tsLayouts {
		for _, variation := range tsFromFileNameVariations(file.Name) {
			t, err := time.Parse(layout, variation)
			if err == nil {
				return t
			}
		}
	}

	// filename could not be parsed by any of the expected patterns
	// so let's default to the created date instead
	return file.CreatedAt
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
