package kyoto

import (
	"net/http"
	"reflect"
)

// PageHandlerFactory is a factory for building Page handler.
// Simple wrapper around RenderPage with context setting.
// Good for defining own project-level handler.
// Example of usage:
// func handle(p ssc.Page) http.HandlerFunc {
//     return func(rw http.ResponseWriter, r *http.Request) {
// 	       ssc.PageHandlerFactory(p, map[string]interface{}{
//	           "internal:rw": rw,
//             "internal:r": r,
//         })(rw, r)
//     }
// }
func PageHandlerFactory(p Page, context map[string]interface{}) http.HandlerFunc {
	// Make page instance
	var pi Page
	pptr := reflect.New(reflect.TypeOf(p).Elem())
	pi = pptr.Interface().(Page)
	// Set context
	for k, v := range context {
		SetContext(pi, k, v)
	}
	// Return handler
	return func(rw http.ResponseWriter, r *http.Request) {
		// Render page
		RenderPage(rw, pi)
		// Clear context
		DelContext(pi, "")
	}
}

// PageHandler is an opinionated net/http handler.
// Context:
// - internal:rw - http.ResponseWritr
// - internal:r - *http.Request
func PageHandler(p Page) http.HandlerFunc {
	// Return handler
	return func(rw http.ResponseWriter, r *http.Request) {
		// Set context

	}
}
