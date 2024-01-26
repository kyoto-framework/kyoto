package component

// Name implements component name getter/setter,
// required for each component type.
type Name struct {
	Name string
}

// SetName is a component name setter.
func (c *Name) SetName(name string) {
	c.Name = name
}

// GetName is a component name getter.
func (c *Name) GetName() string {
	return c.Name
}
