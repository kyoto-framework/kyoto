package kyoto

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStore(t *testing.T) {
	// Initialize store
	store := NewStore()

	// Validate setting works as expected
	store.Set("test", "value")
	if store.state["test"] == nil {
		t.Error("store.Set not working")
	}

	// Validate getting works as expected
	if store.Get("test") != "value" {
		t.Error("store.Get not working")
	}

	// Validate delete works as expected
	store.Del("test")
	if store.state["test"] != nil {
		t.Error("store.Del not working")
	}

	// Validate specific getters
	store.state["internal:rw"] = httptest.NewRecorder()
	store.state["internal:r"], _ = http.NewRequest("GET", "http://localhost:8080/", nil)
	if store.GetResponseWriter() == nil {
		t.Error("store.GetResponseWriter not working")
	}
	if store.GetRequest() == nil {
		t.Error("store.GetRequest not working")
	}
}

func TestStoreNotPanics(t *testing.T) {
	// Initialize store
	store := NewStore()
	// Try to get a non-existent value
	store.Get("test")
}
