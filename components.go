package kyoto

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"runtime"
	"strings"

	"git.sr.ht/~kyoto-framework/zen"
)

// ****************
// Component definition and usage functions
// ****************

// Component represents a kyoto component, defined as a function.
type Component[T any] func(*Context) T

// ComponentF represents a future for a component work result.
// Under the hood it wraps zen.Future[T].
type ComponentF[T any] zen.Future[T]

// MarshalJSON implements future marshalling.
func (f *ComponentF[T]) MarshalJSON() ([]byte, error) {
	return (*zen.Future[T])(f).MarshalJSON()
}

// MarshalJSON implements future unmarshalling.
func (f *ComponentF[T]) UnmarshalJSON(data []byte) error {
	return (*zen.Future[T])(f).UnmarshalJSON(data)
}

// awaitable is an interface to call an await without relying on generics.
type awaitable interface {
	await() any
}

// await is a method to implement awaitable interface
// and utilize zen.Await in a non-generic way.
func (c *ComponentF[T]) await() (val any) {
	val, err := zen.Await((*zen.Future[T])(c))
	if err != nil {
		panic(err)
	}
	return
}

// Use is a function to use your components in your code.
// Triggers component execution and returns a future for a component work result (ComponentF).
//
// Example:
//
//		func CompBar(ctx *kyoto.Context) (state CompBarState) {
//			...
//		}
//
//		func PageFoo(ctx *kyoto.Context) (state PageFooState) {
//			...
//			state.Bar = kyoto.Use(ctx, CompBar)  // Awaitable *kyoto.ComponentF[CompBarState]
//			...
//		}
func Use[T any](c *Context, component Component[T]) *ComponentF[T] {
	return (*ComponentF[T])(zen.Async(func() (T, error) {
		return component(c), nil
	}))
}

// Await accepts any awaitable component and returns it's state.
// It's a function supposed to be used as a template function.
//
// Template example:
//
//		{{ template "CompBar" await .Bar }}
//
// Go example:
//
//		barf = kyoto.Use(ctx, CompBar) // Awaitable *kyoto.ComponentF[CompBarState]
//		...
//		bar = kyoto.Await(barf) // CompBarState
//
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

// ComponentName takes a component function and tries to extract it's name.
// Be careful while using this function, may lead to undefined behavior in case of wrong value.
//
// Example:
//
//		func CompBar(ctx *kyoto.Context) (state CompBarState) {
//			...
//		}
//
//		func main() {
//			fmt.Println(kyoto.ComponentName(CompBar)) // "CompBar"
//		}
//
func ComponentName(component any) string {
	funcpath := runtime.FuncForPC(reflect.ValueOf(component).Pointer()).Name()
	tokens := strings.Split(funcpath, ".")
	if tokens[len(tokens)-1] == "func1" || tokens[len(tokens)-1] == "func2" {
		return tokens[len(tokens)-2]
	} else {
		return tokens[len(tokens)-1]
	}
}

// MarshalState encodes components' state for a client.
// Supposed to be used as a template function.
//
// Template example:
//
//		{{ state . }}
//
// Go example:
//
//		compStateEnc := kyoto.MarshalState(compState)
//
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

// UnmarshalState decodes components' state from a client.
// Supposed to be used internaly only, exposed just in case.
func UnmarshalState(state string, target any) {
	// Decode from base64
	statebts, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		panic("Error while decoding component state. " + err.Error())
	}
	// Deserialize component state from json
	err = json.Unmarshal(statebts, target)
	if err != nil {
		panic("Error while deserializing component state. " + err.Error())
	}
}
