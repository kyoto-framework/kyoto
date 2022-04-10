package bench

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/actions"
	"github.com/kyoto-framework/kyoto/lifecycle"
	"github.com/kyoto-framework/kyoto/render"
)

// componentBench is a usual component with content.
func componentBench(index int) func(*kyoto.Core) {
	return func(c *kyoto.Core) {
		lifecycle.Init(c, func() {
			c.State.Set("internal:name", fmt.Sprintf(`component%v`, index))
			c.State.Set("Content", fmt.Sprintf(`I'm a component with index %v`, index))
		})
	}
}

// componentBenchWriter is a usual component with content,
// but also have a custom rendering.
func componentBenchWriter(index int) func(*kyoto.Core) {
	return func(c *kyoto.Core) {
		lifecycle.Init(c, func() {
			c.State.Set("internal:name", fmt.Sprintf(`component%v`, index))
			c.State.Set("Content", fmt.Sprintf(`I'm a component with index %v`, index))
		})
		render.Writer(c, func(w io.Writer) error {
			_, err := fmt.Fprintf(w, `<div> Content: %s </div>`, c.State.Get("Content"))
			return err
		})
	}
}

// BenchmarkEmpty is testing performance of empty page.
// This bench was created to determine overhead of scheduler architecture.
func BenchmarkEmpty(b *testing.B) {

	// Define a page
	page := func(core *kyoto.Core) {
		render.Template(core, func() *template.Template {
			return template.Must(template.New("bench").Parse(`
				<html>
					<head>
						<title>Kyoto benchmark page</title>
					</head>
					<body>
						I'm a content
					</body>
				</html>
			`))
		})
	}

	// Define a mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", render.PageHandler(page))

	// Bench
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)
	}
}

// BenchmarkComponents is testing performance of page with various components count.
// This bench was created to ensure that rendering time will be acceptable.
func BenchmarkComponents(b *testing.B) {
	counts := []int{1, 100, 1000}
	for _, count := range counts {
		b.Run(fmt.Sprintf("%v", count), func(b *testing.B) {
			// Component definitions and usage
			definitions := ""
			usages := ""
			for i := 0; i < count; i++ {
				definitions += fmt.Sprintf(`{{ define "component%v" }}<div> Content: {{ .Content }} </div>{{ end }}`, i)
				usages += fmt.Sprintf(`{{ template "component%v" .c%v }}`, i, i)
			}

			// Define a page
			page := func(core *kyoto.Core) {
				lifecycle.Init(core, func() {
					for i := 0; i < count; i++ {
						core.Component("c"+strconv.Itoa(i), componentBench(i))
					}
				})
				render.Template(core, func() *template.Template {
					tmpl := ""
					tmpl += definitions
					tmpl += fmt.Sprintf(`
						<html>
							<head>
								<title>Kyoto benchmark page</title>
							</head>
							<body> %s </body>
						</html>
					`, usages)
					return template.Must(template.New("bench").Parse(tmpl))
				})
			}

			// Define a mux
			mux := http.NewServeMux()
			mux.HandleFunc("/", render.PageHandler(page))

			// Bench
			for i := 0; i < b.N; i++ {
				req, _ := http.NewRequest("GET", "/", nil)
				res := httptest.NewRecorder()
				mux.ServeHTTP(res, req)
			}
		})
	}
}

// BenchmarkComponentsDynamicRender is doing the same as BenchmarkComponents,
// but using dynamic rendering with `render` function for classic templates.
// This test was created to measure performance impact of using `render` for classic templates.
func BenchmarkComponentsDynamicRender(b *testing.B) {
	counts := []int{1, 100}
	for _, count := range counts {
		b.Run(fmt.Sprintf("%v", count), func(b *testing.B) {
			// Components definitions and usage
			definitions := ""
			usages := ""
			for i := 0; i < count; i++ {
				definitions += fmt.Sprintf(`{{ define "component%v" }}<div> Content: {{ .Content }} </div>{{ end }}`, i)
				usages += fmt.Sprintf(`{{ render .c%v }}`, i)
			}

			// Define a page
			page := func(core *kyoto.Core) {
				lifecycle.Init(core, func() {
					lifecycle.Init(core, func() {
						for i := 0; i < count; i++ {
							core.Component("c"+strconv.Itoa(i), componentBench(i))
						}
					})
				})
				render.Template(core, func() *template.Template {
					tmpl := ""
					tmpl += definitions
					tmpl += fmt.Sprintf(`
						<html>
							<head>
								<title>Kyoto benchmark page</title>
							</head>
							<body> %s </body>
						</html>
					`, usages)
					return template.Must(template.New("bench").Funcs(render.FuncMap(core)).Parse(tmpl))
				})
			}

			// Define a mux
			mux := http.NewServeMux()
			mux.HandleFunc("/", render.PageHandler(page))

			// Bench
			for i := 0; i < b.N; i++ {
				req, _ := http.NewRequest("GET", "/", nil)
				res := httptest.NewRecorder()
				mux.ServeHTTP(res, req)
			}
		})
	}
}

// BenchmarkComponentsWriter is doing the same as BenchmarkComponents,
// but using dynamic rendering with `render` function for custom rendering.
// This test was created to measure performance impact of using `render` for custom rendering.
func BenchmarkComponentsWriter(b *testing.B) {
	counts := []int{1, 100, 1000}
	for _, count := range counts {
		b.Run(fmt.Sprintf("%v", count), func(b *testing.B) {
			// Components definitions and usage
			usages := ""
			for i := 0; i < count; i++ {
				usages += fmt.Sprintf(`{{ render .c%v }}`, i)
			}
			// Define a page
			page := func(core *kyoto.Core) {
				lifecycle.Init(core, func() {
					lifecycle.Init(core, func() {
						for i := 0; i < count; i++ {
							core.Component("c"+strconv.Itoa(i), componentBenchWriter(i))
						}
					})
				})
				render.Template(core, func() *template.Template {
					tmpl := ""
					tmpl += fmt.Sprintf(`
						<html>
							<head>
								<title>Kyoto benchmark page</title>
							</head>
							<body> %s </body>
						</html>
					`, usages)
					return template.Must(template.New("bench").Funcs(render.FuncMap(core)).Parse(tmpl))
				})
			}

			// Define a mux
			mux := http.NewServeMux()
			mux.HandleFunc("/", render.PageHandler(page))

			// Bench
			for i := 0; i < b.N; i++ {
				req, _ := http.NewRequest("GET", "/", nil)
				res := httptest.NewRecorder()
				mux.ServeHTTP(res, req)
			}
		})
	}
}

// BenchmarkAction is testing performance of actions.
func BenchmarkAction(b *testing.B) {
	// Define a component with an action
	component := func(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.State.Set("Content", "I'm a component content")
		})
		actions.Define(core, "Action", func(args ...interface{}) {
			core.State.Set("Content", "I'm an action content")
		})
	}

	// Define a template builder
	tb := func(c *kyoto.Core) *template.Template {
		return template.Must(template.New("bench").Funcs(render.FuncMap(c)).Parse(`
				{{ define "component" }}
				<div {{ componentattrs . }}>
					{{ .Content }}
				</div>
				{{ end }}
				<html>
					<head>
						<title>Kyoto benchmark page</title>
					</head>
					<body>
						{{ render .c1 }}
					</body>
				</html>
			`))
	}

	// Define a page
	page := func(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.Component("c1", component)
		})
		render.Template(core, func() *template.Template {
			return tb(core)
		})
	}

	// Define a mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", render.PageHandler(page))
	mux.HandleFunc("/internal/actions/", actions.Handler(tb))

	// Register component
	actions.RegisterWithName("component", component)

	// Bench
	for i := 0; i < b.N; i++ {
		reqb := bytes.Buffer{}
		reqw := multipart.NewWriter(&reqb)
		reqw.WriteField("State", `{"Content":"I'm a component content"}`)
		reqw.WriteField("Args", "[]")
		reqw.Close()
		req, _ := http.NewRequest("POST", "/internal/actions/component/Action", &reqb)
		req.Header.Set("Content-Type", reqw.FormDataContentType())
		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)
	}
}
