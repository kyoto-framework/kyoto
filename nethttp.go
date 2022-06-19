package kyoto

import (
	"log"
	"net/http"
)

// ****************
// Page handling
// ****************

// HandlePage registers the page for the given pattern in the DefaultServeMux.
// It's a wrapper around http.HandlePage, but accepts a page instead of usual http.HandlerFunc.
func HandlePage[T any](pattern string, page Component[T]) {
	log.Printf("Registering '%s' page under '%s'", ComponentName(page), pattern)
	http.HandleFunc(pattern, HandlerPage(page))
}

// HandlerPage returns a http.HandlerPage that renders the page.
func HandlerPage[T any](page Component[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize context
		ctx := &Context{
			ResponseWriter: w,
			Request:        r,
		}
		// Trigger building
		state := page(ctx)
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
func Serve(addr string) {
	log.Println("Starting server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
