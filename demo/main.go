package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yuriizinets/go-ssc"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", PageIndexHandler)
	mux.HandleFunc("/SSA/", ssc.SSAHandler)

	if os.Getenv("PORT") == "" {
		log.Println("Listening on localhost:25025")
		http.ListenAndServe("localhost:25025", mux)
	} else {
		log.Println("Listening on 0.0.0.0:" + os.Getenv("PORT"))
		http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	}
}
