package component

import (
	"net/http"
)

// Context is the context of the current request.
// It is passed to the pages and components.
type Context struct {
	// Handler
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	// Store
	Store
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        r,
		Store:          NewMapStore(),
	}
}

type Store interface {
	Get(key string) any
	Set(key string, value any)
}

type MapStore struct {
	store map[string]any
}

func (s *MapStore) Get(key string) any {
	return s.store[key]
}

func (s *MapStore) Set(key string, value any) {
	s.store[key] = value
}

func NewMapStore() *MapStore {
	return &MapStore{
		store: make(map[string]any),
	}
}
