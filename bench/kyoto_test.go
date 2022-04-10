package bench

import (
	"bytes"
	"fmt"
	"html/template"
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

// BenchmarkThousandComponents is testing performance of page with various components count.
// This bench was created to ensure that rendering time will be acceptable.
func BenchmarkComponents(b *testing.B) {
	counts := []int{1, 100, 1000}
	for _, count := range counts {
		b.Run(fmt.Sprintf("%v", count), func(b *testing.B) {
			// Define components with definitions and usage
			components := make([]func(*kyoto.Core), 0)
			definitions := ""
			usages := ""
			for i := 0; i < count; i++ {
				components = append(components, func(core *kyoto.Core) {
					lifecycle.Init(core, func() {
						core.State.Set("Content", "I'm a component content")
					})
				})
				definitions += fmt.Sprintf(`{{ define "component%v" }}<div> Content: {{ .Content }} </div>{{ end }}`, i)
				usages += fmt.Sprintf(`{{ template "component%v" .c%v }}`, i, i)
			}

			// Define a page
			page := func(core *kyoto.Core) {
				lifecycle.Init(core, func() {
					for i, component := range components {
						core.Component("c"+strconv.Itoa(i), component)
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
