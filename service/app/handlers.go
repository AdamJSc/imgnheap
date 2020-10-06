package app

import (
	"imgnheap/service/views"
	"net/http"
)

func indexHandler(c Container) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Templates().ExecuteTemplate(w, "index", views.IndexPage{Page: views.Page{Title: "Select Directory"}})
	}
}
