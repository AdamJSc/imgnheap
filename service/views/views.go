package views

import (
	"bytes"
	"github.com/markbates/pkger"
	"html/template"
	"io"
	"log"
	"os"
)

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

type IndexPage struct {
	Page
}

type Page struct {
	Title string
}
