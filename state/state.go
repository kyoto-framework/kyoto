package state

import (
	"fmt"

	"github.com/kyoto-framework/kyoto"
)

// State is a value wrapper for a comfortable work with kyoto.Core.State
type State[T any] struct {
	Value T
	State bool
}

// Set is a setter for a state
func (s *State[T]) Set(value T) {
	s.Value = value
}

// SetAny is a setter for a state (accepts any)
func (s *State[T]) SetAny(value any) {
	s.Value = value.(T)
}

// Get is a getter for a state
func (s *State[T]) Get() T {
	return s.Value
}

// String is a serializer for a state value
func (s *State[T]) String() string {
	return fmt.Sprintf("%v", s.Value)
}

// New is a constructor for a state
func New[T any](c *kyoto.Core, alias string, value T) *State[T] {
	// Build new state
	s := &State[T]{
		Value: value,
		State: true,
	}
	// Bind to the core state
	c.State.Set(alias, s)
	// Return the state
	return s
}
