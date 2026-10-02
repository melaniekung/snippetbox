package main

import (
	"net/http"
	"runtime/debug"
)


func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var(
		method = r.Method
		uri = r.URL.RequestURI()
		trace = string(debug.Stack())  // convert byte slice to string
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends specific status code and corresponding description to user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}