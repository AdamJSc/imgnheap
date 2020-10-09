package domain

import (
	"imgnheap/service/app"
	"imgnheap/service/models"
)

type ImagesAgentInjector interface {
	app.FileSystemInjector
}

// ImagesAgent encapsulates all of our image-related operations
type ImagesAgent struct {
	ImagesAgentInjector
}

// GetImageFilesFromDirectory returns a slice of the image files present within the provided path
func (i *ImagesAgent) GetImageFilesFromDirectory(path string) ([]models.File, error) {
	files, err := i.FileSystem().GetFilesInDirectory(path)
	if err != nil {
		return nil, err
	}

	imgFileExts := []string{
		"png",
		"jpg",
		"jpeg",
	}

	return filterFilesByExtensions(files, imgFileExts...), nil
}
