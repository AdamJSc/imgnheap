package views

import (
	"bytes"
	"github.com/markbates/pkger"
	"html/template"
	"io"
	"log"
	"os"
)

// MustParseTemplates parses the HTML view templates, otherwise fails on error
func MustParseTemplates() *template.Template {
	tpl := template.New("imgnheap")

	if err := pkger.Walk("/service/views/html", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := pkger.Open(path)
		if err != nil {
			return err
		}
		var b bytes.Buffer
		if _, err = io.Copy(&b, file); err != nil {
			return err
		}
		tpl = template.Must(tpl.Parse(b.String()))
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return tpl
}

// Page represents the dataset required by a standard page template
type Page struct {
	Title          string
	DirPath        string
	WithStartAgain bool
}

// NewPage returns a new Page object
func NewPage(title, dirPath string, withStartAgain bool) Page {
	return Page{
		Title:          title,
		DirPath:        dirPath,
		WithStartAgain: withStartAgain,
	}
}

// IndexPage represents the dataset required by the index page
type IndexPage struct {
	Page
}

// CatalogMethodSelectionPage represents the dataset required by the catalog method selection page
type CatalogMethodSelectionPage struct {
	Page
	ImageFilesCount int
}

// CatalogByTagPage represents the dataset required by the catalog by tag page
type CatalogByTagPage struct {
	Page
	ImageFilesCount int
	ImageFileName   string
	TagsWithCount   map[string]int
}

// ErrorPage represents the dataset required by an error page
type ErrorPage struct {
	Page
	Error struct {
		Code    int
		Message string
		Detail  string
	}
}
