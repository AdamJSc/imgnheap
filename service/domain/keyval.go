package domain

import (
	"fmt"
	"imgnheap/service/app"
)

// InMemoryKeyValStore defines an in-memory key/value store
type InMemoryKeyValStore struct {
	app.KeyValStore
	mem map[string]interface{}
}

// Read implements app.KeyValStore.Read()
func (i *InMemoryKeyValStore) Read(key string) (interface{}, error) {
	val, ok := i.mem[key]
	if !ok {
		return "", NotFoundError{Err: fmt.Errorf("no value found at key %s", key)}
	}

	return val, nil
}

// Write implements app.KeyValStore.Write()
func (i *InMemoryKeyValStore) Write(key string, val interface{}) error {
	i.mem[key] = val
	return nil
}

// NewInMemoryKeyValStore returns a newly-instantiated InMemoryKeyValStore
func NewInMemoryKeyValStore() *InMemoryKeyValStore {
	return &InMemoryKeyValStore{
		mem: make(map[string]interface{}),
	}
}
