package domain

import "fmt"

// InMemoryKeyValStore defines an in-memory key/value store
type InMemoryKeyValStore struct {
	mem map[string]string
}

// Read implements KeyValStore.Read()
func (i *InMemoryKeyValStore) Read(key string) (string, error) {
	val, ok := i.mem[key]
	if !ok {
		return "", NotFoundError{Err: fmt.Errorf("no value found at key %s", key)}
	}

	return val, nil
}

// Write implements KeyValStore.Write()
func (i *InMemoryKeyValStore) Write(key string, val string) error {
	i.mem[key] = val
	return nil
}

// NewInMemoryKeyValStore returns a newly-instantiated InMemoryKeyValStore
func NewInMemoryKeyValStore() *InMemoryKeyValStore {
	return &InMemoryKeyValStore{
		mem: make(map[string]string),
	}
}
