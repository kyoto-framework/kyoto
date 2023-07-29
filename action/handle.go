package action

import (
	"go.kyoto.codes/v3/component"
)

// Handle allows to handle action inside of component.
// Returns handling status (if the action was handled or not).
func Handle(ctx *component.Context, name string, handler func(...any)) bool {
	// Preload context, if not preloaded yet
	context(ctx)
	// Pass, if action not recognized.
	// It means we are not in the action workflow now.
	if ctx.Store.Get("Action") == nil {
		return false
	}
	// Extract action
	action := ctx.Store.Get("Action").(Action)
	// Pass, if action already handled
	if action.handled {
		return false
	}
	// Pass, if action name doesn't match
	if action.Action != name {
		return false
	}
	// Execute handler
	handler(action.Args...)
	// Mark action as handled
	action.handled = true
	// Set action to context
	ctx.Store.Set("Action", action)
	// Return
	return true
}
