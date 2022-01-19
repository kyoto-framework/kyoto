package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// ComponentID is a function to get component id.
// Now it uses pointer address.
func ComponentID(component interface{}) string {
	return fmt.Sprintf("%p", component)
}

// ComponentName is a function to extract component name.
// In case of struct pointer, takes struct name.
// In case of Map, extracts component name with "internal:name" key.
// In case of function, takes function name.
func ComponentName(component interface{}) string {
	// In case of passed function component
	if reflect.TypeOf(component).Kind() == reflect.Func {
		funcpath := runtime.FuncForPC(reflect.ValueOf(component).Pointer()).Name()
		tokens := strings.Split(funcpath, ".")
		if len(tokens) == 1 {
			return tokens[0]
		} else {
			return tokens[1]
		}
	}
	// In case of passed state map (component name must to be stored inside)
	if reflect.TypeOf(component).Kind() == reflect.Map {
		return component.(map[string]interface{})["internal:name"].(string)
	}
	// In case of passed struct pointer
	if reflect.TypeOf(component).Kind() == reflect.Ptr && reflect.TypeOf(component).Elem().Kind() == reflect.Struct {
		return reflect.TypeOf(component).Elem().Name()
	}
	// Default is panic
	panic("Unable to get component name")
}

// ComponentSerialize is a function to serialize component state.
func ComponentSerialize(component interface{}) string {
	// If state is passed, make local cleaned copy
	if component, ok := component.(map[string]interface{}); ok {
		_component := map[string]interface{}{}
		for k, v := range component {
			if !strings.HasPrefix(k, "internal:") {
				_component[k] = v
			}
		}
		component = _component
	}
	// Serialize component state into json
	statebts, err := json.Marshal(component)
	if err != nil {
		println("Error while serializing component state")
		panic(err)
	}
	// Encode to base64 and replace slashes to avoid parsing errors
	state := strings.ReplaceAll(base64.StdEncoding.EncodeToString(statebts), "/", "-")
	// Return
	return state
}
