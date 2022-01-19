package actions

import "github.com/kyoto-framework/kyoto"

// Define is a function to register a new action.
func Define(c *kyoto.Core, name string, action func(args ...interface{})) {
	actions := GetActions(c)
	actions[name] = action
	SetActions(c, actions)
}
