package kyoto

import (
	"encoding/json"
	"fmt"
	"math"
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

func TestAwaitNonAwaitableError(t *testing.T) {
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
			fmt.Println("Recovered in TestAwaitNonAwaitableError", r)
		}
	}()

	// Await non-awaitable component
	uuid := Await(CUUID)

	// Error if not panicked
	t.Errorf("Await returned %s, expected panic", uuid)
}

func TestMarshalState(t *testing.T) {
	// Define component state
	type FooState struct {
		Foo string
	}

	state := FooState{Foo: "Bar"}
	b64state := MarshalState(state)

	if b64state == "" {
		t.Errorf("MarshalState returned empty string, expected non-empty string")
	}

	emptyState := FooState{}
	UnmarshalState(b64state, &emptyState)

	if emptyState.Foo != state.Foo {
		t.Errorf("UnmarshalState returned %s, expected %s", emptyState.Foo, state.Foo)
	}
}

func TestMarshalStateError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestMarshalStateError", r)
		}
	}()

	MarshalState(math.Inf(1))

	t.Error("MarshalState did not panic")
}

func TestUnmarshalDecodingError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestUnmarshalDecodingError", r)
		}
	}()

	UnmarshalState("{\"Foo\":\"Bar\"}", &struct{}{})
	t.Error("UnmarshalState did not panic")
}

func TestUnmarshalDeserializeError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestUnmarshalDeserializeError", r)
		}
	}()

	UnmarshalState("", &struct{}{})
	t.Error("UnmarshalState did not panic")
}
