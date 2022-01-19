package kyoto

import (
	"net/http"
	"sync"
)

// State is an interface around Store than provides guarantees
// for state operations.
type State interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Del(key string)
	Export() map[string]interface{}
}

// Context is an interface around Store than provides guarantees
// for context operations.
type Context interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Del(key string)
	GetResponseWriter() http.ResponseWriter
	GetRequest() *http.Request
}

// Store is a simple wrapper with lock around map[string]interface{}.
type Store struct {
	state map[string]interface{}
	lock  sync.RWMutex
}

// NewStore is a constructor for Store.
func NewStore() *Store {
	return &Store{state: make(map[string]interface{})}
}

// Get is a getter for Store.
func (s *Store) Get(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.state[key]
}

// Set is a setter for Store.
func (s *Store) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.state[key] = value
}

// Del is a delete method for Store.
func (s *Store) Del(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.state, key)
}

// Export is a method to export map[string]interface{} from internals.
func (s *Store) Export() map[string]interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.state
}

// Aliases to avoid verbose type casting

// GetResponseWriter is a getter method with type casting to http.ResponseWriter.
func (s *Store) GetResponseWriter() http.ResponseWriter {
	return s.Get("internal:rw").(http.ResponseWriter)
}

// GetRequest is a getter method with type casting to *http.Request.
func (s *Store) GetRequest() *http.Request {
	return s.Get("internal:r").(*http.Request)
}
