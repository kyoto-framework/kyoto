package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/kyoto-framework/kyoto"
)

func ssatemplate(p kyoto.Page) *template.Template {
	return template.Must(template.New("SSA").Funcs(kyoto.TFuncMap()).ParseGlob("*.html"))
}

func main() {
	// Set flags
	kyoto.INSIGHTS = true
	kyoto.INSIGHTS_CLI = true

	// Init mux
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", kyoto.PageHandler(&PageIndex{}))
	mux.HandleFunc("/raw", kyoto.PageHandler(&PageRaw{}))
	// SSA plugin
	mux.HandleFunc("/SSA/", kyoto.SSAHandler(ssatemplate))

	// Run
	if os.Getenv("PORT") == "" {
		log.Println("Listening on http://localhost:25025")
		http.ListenAndServe("localhost:25025", mux)
	} else {
		log.Println("Listening on http://0.0.0.0:" + os.Getenv("PORT"))
		http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	}
}
