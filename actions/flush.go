package actions

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/kyoto-framework/kyoto"
)

func Flush(b *kyoto.Core) {
	// Extract context
	tb := b.Context.Get("internal:render:tb").(func() *template.Template)
	rw := b.Context.GetResponseWriter()
	rwf := rw.(http.Flusher)
	// Gather state
	for _, job := range b.Scheduler.Jobs {
		if job.Group == "state" {
			job.Func()
		}
	}
	// Render
	buffer := bytes.Buffer{}
	err := tb().Execute(&buffer, b.State.Export())
	if err != nil {
		panic(err)
	}
	html := buffer.String()
	// Remove newlines (not supported by SSA)
	html = strings.ReplaceAll(html, "\n", "")
	// Write SSE
	_, err = fmt.Fprintf(rw, "data: %v\n\n", html)
	if err != nil {
		panic(err)
	}
	// Flush
	rwf.Flush()
}
