package main

import (
	"fmt"
	"html/template"
	"imgnheap/service/app"
	"imgnheap/service/views"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())

	c := container{
		templates: views.MustParseTemplates(),
	}

	port := 8080
	router := app.RegisterRouter(c)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("listening on port %d...\n", port)
	log.Fatal(server.ListenAndServe())
}

type container struct {
	templates *template.Template
}

func (c container) Templates() *template.Template {
	return c.templates
}
