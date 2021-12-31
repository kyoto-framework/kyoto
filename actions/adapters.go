package actions

import "github.com/kyoto-framework/kyoto"

func Define(c *kyoto.Core, name string, action func(args ...interface{})) {
	actions := GetActions(c)
	actions[name] = action
	SetActions(c, actions)
}
