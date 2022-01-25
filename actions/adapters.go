package actions

import "github.com/kyoto-framework/kyoto"

// Define is a function to define a new action for a component.
func Define(c *kyoto.Core, name string, action func(args ...interface{})) {
	actions := GetActions(c)
	actions[name] = action
	SetActions(c, actions)
}
