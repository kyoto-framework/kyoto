package component

// State is a component state.
// State implementation may vary depending on the component type.
// State also holds optional/additional parameters and information,
// like rendering type/strategy, component name, etc.
type State interface {
	// Marshal state into string representation.
	// It's required for injecting state into DOM.
	Marshal(state any) string
	// Unmarshal state from string representation.
	Unmarshal(str string, state any) error

	// We can't extract component name without
	// having original Component function,
	// so we are giving an ability to inject component name
	// directly into resulting State.
	// Name is required for operations like rendering.

	// GetName is a component name getter.
	GetName() string
	// SetName is a component name setter.
	SetName(name string)
}
