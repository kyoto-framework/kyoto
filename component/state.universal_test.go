package component_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.kyoto.codes/v3/component"
)

func TestMarshal(t *testing.T) {
	t.Parallel()
	// Define initial states
	stateOriginal := component.Use(component.NewContext(nil, nil), testComponent)()
	state := component.Use(component.NewContext(nil, nil), testComponent)()
	// Marshalling state to string
	stateMarshalled := state.Marshal()
	// Unmarshalling state from string
	state.Unmarshal(stateMarshalled)
	// Assert results
	assert.True(t, assert.ObjectsAreEqualValues(stateOriginal, state), "Something wrong with state.Universal marshalling")
}
