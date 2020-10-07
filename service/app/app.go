package app

import (
	"html/template"
)

// Container defines our app's container interface
type Container interface {
	TemplatesInjector
	KeyValStoreInjector
}

type TemplatesInjector interface{ Templates() *template.Template }
type KeyValStoreInjector interface{ KeyValStore() KeyValStore }

// KeyValStore defines operations for transacting with key/value storage
type KeyValStore interface {
	Read(key string) (string, error)
	Write(key string, val string) error
}
