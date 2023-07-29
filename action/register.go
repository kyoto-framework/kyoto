package action

import "go.kyoto.codes/v3/component"

var components = map[string]component.Component{}

// Register registers a component into the library global store
// to be able to find and render specific component during action execution.
func Register(c component.Component) {
	components[c.GetName()] = c
}
