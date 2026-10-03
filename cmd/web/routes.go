package main

import "net/http"

// returns servemux containing application routes
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// create file server to serve files
	// path is relative to project root directory
	fileServer := http.FileServer(http.Dir("./assets/static"))

	// register file server as handler for all URL paths starting with "/static/"
	// strip "/static" prefix before request reaches file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// register handler functions and route patterns
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return mux
}
