package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/yuriizinets/kyoto"
)

func ssatemplate(p kyoto.Page) *template.Template {
	return template.Must(template.New("SSA").Funcs(tfuncs()).ParseGlob("*.html"))
}

func ssahandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		kyoto.SSAHandlerFactory(ssatemplate, map[string]interface{}{
			"internal:rw": rw,
			"internal:r":  r,
		})(rw, r)
	}
}

func main() {
	mux := http.NewServeMux()

	// Statics
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/dist"))))

	// Routes
	mux.HandleFunc("/", kyoto.PageHandler(&PageIndex{}))
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/.vuepress/dist"))))
	// SSA plugin
	mux.HandleFunc("/SSA/", ssahandler())

	// Run
	if os.Getenv("PORT") == "" {
		log.Println("Listening on localhost:25025")
		http.ListenAndServe("localhost:25025", mux)
	} else {
		log.Println("Listening on 0.0.0.0:" + os.Getenv("PORT"))
		http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	}
}
