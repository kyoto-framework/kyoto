package actions

import (
	"sync"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/helpers"
)

var registry = map[string]func(*kyoto.Core){}
var registryrw = sync.RWMutex{}

// Register is a function to register a dynamic component.
func Register(components ...func(*kyoto.Core)) {
	// Acquire write lock
	registryrw.Lock()
	defer registryrw.Unlock()
	// Iterate and register components
	for _, component := range components {
		registry[helpers.ComponentName(component)] = component
	}

}

// RegisterWithName is a function to register a dynamic component with a specified name.
func RegisterWithName(name string, component func(*kyoto.Core)) {
	// Acquire write lock
	registryrw.Lock()
	defer registryrw.Unlock()
	// Register component
	registry[name] = component
}
