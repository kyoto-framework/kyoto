package component

import (
	"go.kyoto.codes/zen/v3/async"
	"go.kyoto.codes/zen/v3/errorsx"
)

// Use allows you to use your components in asynchronous way.
// It's a basic and preferred way to use your components.
func Use(ctx *Context, component Component) Future {
	// Create state future.
	ftr := async.New(func() (State, error) {
		return component(ctx), nil
	})
	// Create and return getter.
	return func() State {
		// Await for state.
		state := errorsx.Must(ftr.Await())
		// Set component name.
		state.SetName(component.GetName())
		// Return state.
		return state
	}
}
