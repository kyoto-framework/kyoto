package component_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.kyoto.codes/v3/component"
)

// Disposable component.
func testDisposableComponent(_ *component.Context) component.State {
	// Set state
	state := &testComponentDisposableState{}
	state.Value = "Default value"
	// Return
	return state
}

func TestMarshalDisposable(t *testing.T) {
	t.Parallel()
	// Define initial states
	stateOriginal := component.Use(component.NewContext(nil, nil), testDisposableComponent)()
	state := component.Use(component.NewContext(nil, nil), testDisposableComponent)()
	// Marshalling state to string
	stateMarshalled := state.Marshal()
	// Assert disposable marshalled state
	assert.Equal(t, "disposable", stateMarshalled)
	// Unmarshalling state from string
	state.Unmarshal(stateMarshalled)
	// Assert results
	assert.True(t, assert.ObjectsAreEqualValues(stateOriginal, state), "Something wrong with state.Disposable marshalling")
}
