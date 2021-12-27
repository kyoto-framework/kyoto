package kyoto

import "net/http"

func PageHandler(page interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Adapt behavior, depending on the page type
		switch page := page.(type) {
		case func(*Builder): // Builder receiver
			b := NewBuilder()
			b.Context.Set("internal:rw", rw)
			b.Context.Set("internal:r", r)
			page(b)
			b.Render(rw)
			b.Execute()
		default: // Not supported
			panic("Page type is not supported")
		}
	}
}
