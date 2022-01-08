package actions

import (
	"sync"

	"github.com/kyoto-framework/kyoto/helpers"
)

var registry = map[string]interface{}{}
var registryrw = sync.RWMutex{}

func Register(components ...interface{}) {
	// Acquire write lock
	registryrw.Lock()
	defer registryrw.Unlock()
	// Iterate and register components
	for _, component := range components {
		registry[helpers.ComponentName(component)] = component
	}

}

func RegisterWithName(name string, component interface{}) {
	// Acquire write lock
	registryrw.Lock()
	defer registryrw.Unlock()
	// Register component
	registry[name] = component
}
