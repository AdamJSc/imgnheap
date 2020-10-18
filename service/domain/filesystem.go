package domain

import (
	"imgnheap/service/app"
	"imgnheap/service/models"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	SubDirByDate = "by-date"
	SubDirByTag  = "by-tag"
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

	dirPath = path.Clean(dirPath)

	if err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			// an error has already occurred
			return err
		}
		if path.Dir(filePath) != dirPath {
			// we're beyond the first directory level, so return early
			return nil
		}
		if info.IsDir() {
			// we're only interested in files, so return early
			return nil
		}

		fileName, ext := ParseNameAndExtensionFromFileName(info.Name())
		modTime := info.ModTime()
		files = append(files, models.NewFile(fileName, ext, dirPath, &modTime))

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

// GetDirectoriesInDirectory implements app.FileSystem.GetDirectoriesInDirectory()
func (o *OsFileSystem) GetDirectoriesInDirectory(dirPath string) ([]models.Directory, error) {
	var dirs []models.Directory

	dirPath = path.Clean(dirPath)

	if err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			// an error has already occurred
			return err
		}
		if path.Dir(filePath) != dirPath {
			// we're beyond the first directory level, so return early
			return nil
		}
		if !info.IsDir() {
			// we're only interested in directories, so return early
			return nil
		}

		dir := models.Directory{
			Name:    info.Name(),
			DirPath: dirPath,
		}

		dirs = append(dirs, dir)

		return nil
	}); err != nil {
		return nil, err
	}

	return dirs, nil
}

// GetContents implements app.FileSystem.GetContents()
func (o *OsFileSystem) GetContents(file models.File) ([]byte, error) {
	contents, err := ioutil.ReadFile(file.FullPath())

	if err != nil {
		// presume this means the file can't be found for all intents and purposes
		// for granular control, look to determine the exact nature of the error message first
		return nil, NotFoundError{Err: err}
	}

	return contents, nil
}

// Copy implements app.FileSystem.Copy()
func (o *OsFileSystem) Copy(file models.File, destDir string) error {
	src, err := os.Open(file.FullPath())
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	destPath := path.Join(destDir, file.NameWithExt())
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	return err
}

// Move implements app.FileSystem.Move()
func (o *OsFileSystem) Move(file models.File, destDir string) error {
	// copy file
	if err := o.Copy(file, destDir); err != nil {
		return err
	}

	// delete original
	if err := os.Remove(file.FullPath()); err != nil {
		return err
	}

	return nil
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
func (f *FileSystemAgent) GetFilesFromDirectoryByExtension(dir string, exts ...string) ([]models.File, error) {
	files, err := f.FileSystem().GetFilesInDirectory(dir)
	if err != nil {
		return nil, err
	}

	if len(exts) == 0 {
		// no filtering required
		return files, nil
	}

	var filtered []models.File

	for _, file := range files {
		if contains(exts, file.Ext) {
			filtered = append(filtered, file)
		}
	}

	return filtered, nil
}

// ProcessFileByCopy copies the provided file to the provided destination directory
func (f *FileSystemAgent) ProcessFileByCopy(file models.File, destDir string) error {
	if err := f.FileSystem().Copy(file, destDir); err != nil {
		return err
	}
	return nil
}

// ProcessFileByMove moves the provided file to the provided destination directory
func (f *FileSystemAgent) ProcessFileByMove(file models.File, destDir string) error {
	if err := f.FileSystem().Move(file, destDir); err != nil {
		return err
	}
	return nil
}

// Stream writes the contents of the provided file to the provided response writer
func (f *FileSystemAgent) Stream(file models.File, w http.ResponseWriter) error {
	contents, err := f.FileSystem().GetContents(file)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", http.DetectContentType(contents))
	w.Header().Set("Content-Length", strconv.Itoa(len(contents)))
	w.Write(contents)

	return nil
}

// GetDirectoriesWithFileCountByExtension returns a slice of the directories present within the provided directory path
// including the count of files within each one that has one of the provided extensions
func (f *FileSystemAgent) GetDirectoriesWithFileCountByExtension(dir string, exts ...string) ([]models.Directory, error) {
	dirs, err := f.FileSystem().GetDirectoriesInDirectory(dir)
	if err != nil {
		return nil, err
	}

	for idx := range dirs {
		files, err := f.GetFilesFromDirectoryByExtension(dirs[idx].FullPath(), exts...)
		if err != nil {
			return nil, err
		}

		dirs[idx].FileCount = len(files)
	}

	return dirs, nil
}

// ParseNameAndExtensionFromFileName returns the name and extension from the provided filename string
func ParseNameAndExtensionFromFileName(fileName string) (string, string) {
	var ext string
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

// GetDestinationDirByDate returns a directory path based on the timestamp parsed from the provided file
func GetDestinationDirByDate(file models.File, sess *models.Session) string {
	if sess == nil {
		return ""
	}

	return path.Join(sess.FullDir(SubDirByDate), file.Ext, ParseTimestampFromFile(file).Format("2006-01-02"))
}

// GetDestinationDirByTag returns a directory path based on the provided session and tag
func GetDestinationDirByTag(sess *models.Session, tag string) string {
	if sess == nil {
		return ""
	}

	return path.Join(sess.FullDir(SubDirByTag), tag)
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
