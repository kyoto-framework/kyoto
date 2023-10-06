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

// Initialize a new context, that will be passed through the components.
// Uses MapStore as a store by default.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        r,
		Store:          NewMapStore(),
	}
}

// Store allows you to store own data inside of the context.
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
