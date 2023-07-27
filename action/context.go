package action

import (
	"github.com/kyoto-framework/kyoto/v3/component"
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
