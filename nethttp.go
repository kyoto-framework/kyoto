package kyoto

import (
	"log"
	"net/http"
	"reflect"
	"time"
)

// ****************
// Page handling
// ****************

// HandlePage registers the page for the given pattern in the DefaultServeMux.
// It's a wrapper around http.HandlePage, but accepts a page instead of usual http.HandlerFunc.
//
// Example:
//
//		func PageFoo(ctx *kyoto.Context) (state PageFooState) { ... }
//
//		func main() {
//			...
//			kyoto.HandlePage("/foo", PageFoo)
//			...
//		}
//
func HandlePage[T any](pattern string, page Component[T]) {
	log.Printf("Registering page '%s':\t'%s'", ComponentName(page), pattern)
	http.HandleFunc(pattern, HandlerPage(page))
}

// HandlerPage returns a http.HandlerPage that renders the page.
//
// Example:
//
//		func PageFoo(ctx *kyoto.Context) (state PageFooState) { ... }
//
//		func main() {
//			...
//			http.HandleFunc("/foo", kyoto.HandlerPage(PageFoo))
//			...
//		}
//
func HandlerPage[T any](page Component[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize context
		ctx := &Context{
			ResponseWriter: w,
			Request:        r,
		}
		// Trigger building
		measure := time.Now()
		state := page(ctx)
		logf("page\t%v: %v", reflect.TypeOf(page), time.Since(measure))
		// Render
		if err := ctx.Template.Execute(w, state); err != nil {
			panic(err)
		}
	}
}

// ****************
// Additional utilities
// ****************

// Serve is a simple wrapper around http.ListenAndServe,
// which will log server starting and will panic on error.
//
// Example:
//
//		func main() {
//			...
//			kyoto.Serve(":8000")
//		}
//
func Serve(addr string) {
	log.Println("Starting server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
