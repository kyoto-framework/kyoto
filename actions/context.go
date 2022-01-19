package actions

import "github.com/kyoto-framework/kyoto"

// GetActions is a function that extracts actions map from context.
func GetActions(b *kyoto.Core) map[string]func(...interface{}) {
	if b.Context.Get("internal:actions") == nil {
		return map[string]func(...interface{}){}
	}
	return b.Context.Get("internal:actions").(map[string]func(...interface{}))
}

// SetActions is a function that injects actions map to context.
func SetActions(b *kyoto.Core, actions map[string]func(...interface{})) {
	b.Context.Set("internal:actions", actions)
}
