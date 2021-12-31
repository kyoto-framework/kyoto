package actions

import "github.com/kyoto-framework/kyoto"

func GetActions(b *kyoto.Core) map[string]func(...interface{}) {
	if b.Context.Get("internal:actions") == nil {
		return map[string]func(...interface{}){}
	}
	return b.Context.Get("internal:actions").(map[string]func(...interface{}))
}

func SetActions(b *kyoto.Core, actions map[string]func(...interface{})) {
	b.Context.Set("internal:actions", actions)
}
