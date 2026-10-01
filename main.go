package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// define a home handler function
// writes a byte slice as response body
// displays home page
func home(w http.ResponseWriter, r *http.Request) {
	// add a 'Server: Go' header to the response header map
	w.Header().Add("Server", "Go")

	w.Write([]byte("Hello from Snippetbox"))
}

// snippetView handler function
// display a specific snippet
func snippetView(w http.ResponseWriter, r *http.Request) {
	// extract value of the id wildcard from request
	// convert to int
	// id must be a positive int
	id, err := strconv.Atoi(r.PathValue(("id")))
	if err != nil && id < 1 {
		http.NotFound(w, r)
		return
	}

	// interpolate id value with message
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// snippetCreate handler function
// display a form for creating a new snippet
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// snippetCreatePost handler function
// save new a snippet
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// send a 201 created status code
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Save a new snippet..."))
}

func main() {
	// initialize new server using http.NewServeMux()
	// register handler functions
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :4000")

	// start new web server using http.ListenAndServe()
	// passing TCP network address (:4000) and servemux
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
