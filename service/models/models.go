package models

import "time"

// File represents a single file
type File struct {
	Name      string
	Ext       string
	CreatedAt time.Time
}
