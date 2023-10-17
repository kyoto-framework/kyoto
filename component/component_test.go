package component_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.kyoto.codes/v3/component"
)

func testComponentWrapped(value string) component.Component {
	return func(_ *component.Context) component.State {
		// Set state
		state := &testComponentUniversalState{}
		state.Value = value
		// return
		return state
	}
}

// TestComponentName ensures ComponentName returns correct values for expected arguments.
func TestComponentName(t *testing.T) {
	t.Parallel()
	// Classic component
	name := component.Use(component.NewContext(nil, nil), testComponent)().GetName()
	assert.Equal(t, "testComponent", name)
	// Wrapped component
	name = component.Use(component.NewContext(nil, nil), testComponentWrapped("Custom Value"))().GetName()
	assert.Equal(t, "testComponentWrapped", name)
}
