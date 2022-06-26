package kyoto

import (
	"fmt"
	"math"
	"testing"
)

// Common test items

type testComponentState struct {
	Value string
}

func testComponent(ctx *Context) (state testComponentState) {
	// Set state
	state.Value = "Default value"
	// Return
	return
}

func testComponentWrapped(value string) Component[testComponentState] {
	return func(ctx *Context) (state testComponentState) {
		// Set state
		state.Value = value
		// return
		return
	}
}

// TestComponentName ensures ComponentName returns correct values for expected arguments.
func TestComponentName(t *testing.T) {
	// Classic component
	name := ComponentName(testComponent)
	if name != "testComponent" {
		t.Errorf("ComponentName returned %s, expected %s", name, "testComponent")
	}
	// Wrapped component
	name = ComponentName(testComponentWrapped)
	if name != "testComponentWrapped" {
		t.Errorf("ComponentName returned %s, expected %s", name, "testComponentWrapped")
	}
}

// TestComponentUse ensures Use returns awaitable future.
func TestComponentUse(t *testing.T) {
	// Create context
	c := &Context{}
	// Use component
	var stateftr any = Use(c, testComponent)
	// Check it's an awaitable future
	if _, implements := stateftr.(awaitable); !implements {
		t.Error("Use is not returning awaitable future")
	}
}

// TestComponentAwait ensures Await returns expected result on awaitable future.
func TestComponentAwait(t *testing.T) {
	// Create context
	c := &Context{}
	// Use component
	var stateftr awaitable = Use(c, testComponent)
	// Check value
	expected := testComponentState{Value: "Default value"}
	if Await(stateftr) != expected {
		t.Error("State from a future is not as expected")
	}
}

// TestComponentAwaitErr ensures Await will panic when passing non-awaitable object.
func TestComponentAwaitErr(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestComponentAwaitErr", r)
		}
	}()
	// Trigger a panic with function awaiting, instead of future awaiting
	val := Await(testComponent)
	// Error if not panicked
	t.Errorf("Await returned %s, expected panic", val)
}

// TestComponentMarshalState ensures component state (un)marshalling consistent.
func TestComponentMarshalState(t *testing.T) {
	// Define a test state
	state := testComponentState{
		Value: "Default value",
	}
	// Marshal
	stateenc := MarshalState(state)
	// Unmarshal
	statedec := testComponentState{}
	UnmarshalState(stateenc, &statedec)
	// Assert
	if statedec != state {
		t.Error("Something wrong with state marshalling")
	}
}

// TestComponentMarshalStateErr ensures state marshalling will panic on incorrect value.
func TestComponentMarshalStateErr(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestComponentMarshalStateErr", r)
		}
	}()
	// Trigger panic
	MarshalState(math.Inf(1))
	// Panic expected
	t.Error("Expected panic")
}

// TestComponentUnmarshalStateErr ensures state unmarshalling will panic on incorrect value.
func TestComponentUnmarshalStateErr(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestComponentUnmarshalStateErr", r)
		}
	}()
	// Trigger panic
	UnmarshalState("{\"Foo\":\"Bar\"}", &struct{}{})
	// Panic expected
	t.Error("Expected panic")
}
