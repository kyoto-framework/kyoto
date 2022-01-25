package actions

import (
	"fmt"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/helpers"
)

type ActionMap map[string]func(args ...interface{})

// GetActions is a function that extracts actions map from context.
func GetActions(c *kyoto.Core) ActionMap {
	contextaddr := fmt.Sprintf("internal:actions:%s", helpers.ComponentID(c))
	if actions, ok := c.Context.Get(contextaddr).(ActionMap); ok {
		return actions
	} else {
		return make(ActionMap)
	}
}

// SetActions is a function that injects actions map to context.
// Avoid using this function directly, Define function is a preffered way.
func SetActions(c *kyoto.Core, actions ActionMap) {
	contextaddr := fmt.Sprintf("internal:actions:%s", helpers.ComponentID(c))
	c.Context.Set(contextaddr, actions)
}
