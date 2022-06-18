/*
	-

	Components

	Kyoto provides a way to define components.
	It's a very common approach for modern libraries to manage frontend parts.
	In kyoto each component is a context receiver, which returns it's state.
	Each component becomes a part of the page or top-level component,
	which executes component asynchronously and gets a state future object.
	In that way your components are executing in a non-blocking way.
*/
package kyoto

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"runtime"
	"strings"

	"github.com/kyoto-framework/zen"
)

// ****************
// Component definition and usage functions
// ****************

// Component represents a kyoto component, defined as a function.
type Component[T any] func(*Context) T

// ComponentF represents a future for a component work result.
// Under the hood it wraps zen.Future[T].
type ComponentF[T any] zen.Future[T]

// awaitable is an interface to call an await without relying on generics.
type awaitable interface {
	await() any
}

// await is a method to utilize zen.Await in a non-generic way.
func (c ComponentF[T]) await() (val any) {
	val, err := zen.Await(zen.Future[T](c))
	if err != nil {
		panic(err)
	}
	return
}

// Use is a function to use your components in your code.
// Triggers component execution and returns a future for a component work result (ComponentF).
func Use[T any](c *Context, component Component[T]) ComponentF[T] {
	return ComponentF[T](zen.Async(func() (T, error) {
		return component(c), nil
	}))
}

// Await accepts any awaitable component and returns it's state.
// It's a function supposed to be used as a template function.
func Await(component any) any {
	if component, implements := component.(awaitable); implements {
		return component.await()
	} else {
		panic("calling await for a non-awaitable object")
	}
}

// ****************
// Component utilities
// ****************

func ComponentName[T any](component Component[T]) string {
	funcpath := runtime.FuncForPC(reflect.ValueOf(component).Pointer()).Name()
	tokens := strings.Split(funcpath, ".")
	if tokens[len(tokens)-1] == "func1" || tokens[len(tokens)-1] == "func2" {
		return tokens[len(tokens)-2]
	} else {
		return tokens[len(tokens)-1]
	}
}

func MarshalState(state any) string {
	// Serialize component state into json
	statebts, err := json.Marshal(state)
	if err != nil {
		panic("Error while marshaling component state: " + err.Error())
	}
	// Encode to base64
	stateb64 := base64.StdEncoding.EncodeToString(statebts)
	// Return
	return stateb64
}

func UnmarshalState(state string, target any) {
	// Deserialize component state from json
	err := json.Unmarshal([]byte(state), target)
	if err != nil {
		panic("Error while unmarshaling component state: " + err.Error())
	}
}

// Uncomment when client will be migrated
// func UnmarshalState(state string, target any) {
// 	// Decode from base64
// 	statebts, err := base64.StdEncoding.DecodeString(state)
// 	if err != nil {
// 		panic("Error while decoding component state. " + err.Error())
// 	}
// 	// Deserialize component state from json
// 	err = json.Unmarshal(statebts, target)
// 	if err != nil {
// 		panic("Error while deserializing component state. " + err.Error())
// 	}
// }
