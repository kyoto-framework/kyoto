package action

import (
	"errors"

	"github.com/kyoto-framework/kyoto/v3/component"
)

// State allows to extract component state from action.
// It's an important part of the action workflow that
// might be easily missed and have to be mentioned explicitly in the docs.
func State(ctx *component.Context, state component.State) error {
	// Prepare context
	context(ctx)
	// Pass, if action not recognized.
	// It means we are not in the action workflow now.
	if ctx.Store.Get("Action") == nil {
		return errors.New("action not recognized")
	}
	// Extract action
	action := ctx.Store.Get("Action").(Action)
	// Pass, if action not recognized
	if action.Component == "" {
		return nil
	}
	// Unmarshal state
	return state.Unmarshal(action.State, state)
}
