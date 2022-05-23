package state

import (
	"encoding/json"
	"fmt"
	"reflect"

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
	// Final value, that will be assigned to the state
	var _value T
	// Act in a different way, depending on the value type
	if reflect.TypeOf(value) == reflect.TypeOf(_value) { // If value is the same type as the state, just set it
		_value = value.(T)
	} else if reflect.TypeOf(value).Kind() == reflect.Map { // If value is a map (f.e. on action), try to map it to the state type
		bts, _ := json.Marshal(value)
		json.Unmarshal(bts, &_value)
	}
	// Set the final value
	s.Value = _value
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
