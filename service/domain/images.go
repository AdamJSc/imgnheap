package domain

import (
	"imgnheap/service/app"
	"imgnheap/service/models"
	"strings"
	"time"
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
		"mp4",
	}

	return filterFilesByExtensions(files, imgFileExts...), nil
}

// ParseTimestampFromFile attempts to parse a timestamp from the provided file
func (i *ImagesAgent) ParseTimestampFromFile(file models.File) time.Time {
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
