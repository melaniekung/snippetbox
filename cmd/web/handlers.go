package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// add header to response header map
	// header name: "Server", header value: "Go"
	w.Header().Add("Server", "Go")

	// initialize slice for paths
	// NOTE: base template must be first file
	files := []string{
		"./assets/html/base.html",
		"./assets/html/partials/nav.html",
		"./assets/html/pages/home.html",
	}

	// read template file into template set
	// NOTE: use ... to pass contents of slice as variadic arguments
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// write template content as response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// extract value of wildcard from request
	// try to convert to integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}
