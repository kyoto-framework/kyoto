package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/kyoto-framework/kyoto/actions"
	"github.com/kyoto-framework/kyoto/render"
)

func setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", render.PageHandler(PageIndex))
}

func setupActions(mux *http.ServeMux) {
	// Register Actions handler
	mux.HandleFunc("/SSA/", actions.Handler(func() *template.Template {
		return template.Must(template.New("Actions").Funcs(render.FuncMap()).ParseGlob("*.html"))
	}))
	// Register Actions components
	actions.Register(
		ComponentUUID(""),
	)
}

func main() {
	// Init mux
	mux := http.NewServeMux()

	// Setup
	setupActions(mux)
	setupRoutes(mux)

	// Run
	if os.Getenv("PORT") == "" {
		log.Println("Listening on http://localhost:25025")
		http.ListenAndServe("localhost:25025", mux)
	} else {
		log.Println("Listening on http://0.0.0.0:" + os.Getenv("PORT"))
		http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	}
}
