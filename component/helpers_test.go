package component_test

import "go.kyoto.codes/v3/component"

// Common test items.
type testComponentUniversalState struct {
	component.Universal
	Value string
}

type testComponentDisposableState struct {
	component.Disposable
	Value string
}

type testComponentServerState struct {
	component.Server
	Value string
}

func testComponent(_ *component.Context) component.State {
	// Set state
	state := &testComponentUniversalState{}
	state.Value = "Default value"
	// Return
	return state
}
