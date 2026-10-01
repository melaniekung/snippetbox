package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// create file server to serve files
	// path is relative to project root directory
	fileServer := http.FileServer(http.Dir("./assets/static/"))

	// register file server as handler for all URL paths starting with "/static/"
	// strip "/static" prefix before request reaches file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// register handler functions and route patterns
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :4000")

	// start new web server
	// params: TCP network address, servermux
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
