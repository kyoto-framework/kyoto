package action

import (
	"go.kyoto.codes/v3/component"
)

// context extracts action parameters
// from request inside of component.Context and writes
// Action to the component.Context store.
func context(ctx *component.Context) {
	// Pass, if already unmarshalled
	if ctx.Store.Get("Action") != nil {
		return
	}
	// Try to unmarshal action parameters
	action := Action{}
	err := action.UnmarshalHttpRequest(ctx.Request)
	if err != nil {
		return
	}
	// Set action parameters to context
	ctx.Store.Set("Action", action)
}
