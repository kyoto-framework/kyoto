package state

import (
	"fmt"

	"github.com/kyoto-framework/kyoto"
)

type State[T any] struct {
	Value T
	State bool
}

func (s *State[T]) Set(value T) {
	s.Value = value
}

func (s *State[T]) SetAny(value any) {
	s.Value = value.(T)
}

func (s *State[T]) Get() T {
	return s.Value
}

func (s *State[T]) String() string {
	return fmt.Sprintf("%v", s.Value)
}

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
