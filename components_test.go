package kyoto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type CUUIDStateTestComponent struct {
	UUID string
}

func CUUIDTestComponent(ctx *Context) (state CUUIDStateTestComponent) {
	// Fetch uuid data
	resp, _ := http.Get("http://httpbin.org/uuid")
	data := map[string]string{}
	json.NewDecoder(resp.Body).Decode(&data)
	// Set state
	state.UUID = data["uuid"]
	return state
}

func TestComponent(t *testing.T) {
	// Create context
	c := Context{}
	// Use component
	Use(&c, CUUIDTestComponent)
	name := ComponentName(CUUIDTestComponent)
	if name != "CUUIDTestComponent" {
		t.Errorf("ComponentName returned %s, expected %s", name, "CUUID")
	}

	type CUUIDState struct {
		UUID string
	}

	uuid := func(ctx *Context) (state CUUIDState) {
		// Fetch uuid data
		resp, _ := http.Get("http://httpbin.org/uuid")
		data := map[string]string{}
		json.NewDecoder(resp.Body).Decode(&data)
		// Set state
		state.UUID = data["uuid"]
		return state
	}

	name = ComponentName(uuid)
	if name != "TestComponent" {
		t.Errorf("ComponentName returned %s, expected %s", name, "TestComponent")
	}
}

func TestAwait(t *testing.T) {
	// Define component state
	type CUUIDState struct {
		UUID string
	}

	// Define component
	CUUID := func(ctx *Context) (state CUUIDState) {
		// Fetch uuid data
		resp, _ := http.Get("http://httpbin.org/uuid")
		data := map[string]string{}
		json.NewDecoder(resp.Body).Decode(&data)
		// Set state
		state.UUID = data["uuid"]
		return state
	}

	// Create context
	c := Context{}

	// Use component
	cuuid := Use(&c, CUUID)

	// Await component
	uuid := Await(cuuid)

	// Check state
	if uuid.(CUUIDState).UUID == "" {
		t.Errorf("Await returned %s, expected non-empty string", uuid)
	}
}

func TestAwaitError(t *testing.T) {
	// Define component state
	type CUUIDState struct {
		UUID string
	}

	// Define component
	CUUID := func(ctx *Context) (state CUUIDState) {
		// Fetch uuid data
		resp, _ := http.Get("http://httpbin.org/uuid")
		data := map[string]string{}
		json.NewDecoder(resp.Body).Decode(&data)
		// Set state
		state.UUID = data["uuid"]
		return state
	}

	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestAwait", r)
		}
	}()

	// Await non-awaitable component
	uuid := Await(CUUID)

	// Error if not panicked
	t.Errorf("Await returned %s, expected panic", uuid)
}
