package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/yuriizinets/go-ssc"
)

func ssatemplate() *template.Template {
	return template.Must(template.New("SSA").Funcs(ssc.Funcs()).ParseGlob("*.html"))
}

func ssahandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ssc.SSAHandlerFactory(ssatemplate, map[string]interface{}{
			"internal:rw": rw,
			"internal:r":  r,
		})(rw, r)
	}
}

func pagehandler(p ssc.Page) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ssc.PageHandlerFactory(p, map[string]interface{}{
			"internal:rw": rw,
			"internal:r":  r,
		})(rw, r)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", pagehandler(&PageIndex{}))
	mux.HandleFunc("/SSA/", ssahandler())

	if os.Getenv("PORT") == "" {
		log.Println("Listening on localhost:25025")
		http.ListenAndServe("localhost:25025", mux)
	} else {
		log.Println("Listening on 0.0.0.0:" + os.Getenv("PORT"))
		http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	}
}
