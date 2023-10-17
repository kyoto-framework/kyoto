package component_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.kyoto.codes/v3/component"
)

func testServerComponent(_ *component.Context) component.State {
	// Set state
	state := &testComponentServerState{}
	state.Value = "Default value"
	state.Timeout = time.Millisecond * 500
	// Return
	return state
}

func TestMarshalServer(t *testing.T) {
	t.Parallel()
	// Define initial states
	stateOriginal := component.Use(component.NewContext(nil, nil), testServerComponent)()
	state := component.Use(component.NewContext(nil, nil), testServerComponent)()
	// Marshalling state to string
	stateMarshalled := state.Marshal()
	// Wait for expiring state file
	time.Sleep(time.Second * 1)
	// Unmarshalling state from string
	state.Unmarshal(stateMarshalled)
	// Assert results
	assert.True(t, assert.ObjectsAreEqualValues(stateOriginal, state), "Something wrong with state.Server marshalling")
}
