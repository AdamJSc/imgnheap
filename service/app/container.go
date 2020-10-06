package app

import "html/template"

type Container interface {
	TemplateInjector
}

type TemplateInjector interface{ Templates() *template.Template }
