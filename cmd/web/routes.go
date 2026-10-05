package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// create file server to serve static files (path relative to root)
	fileServer := http.FileServer(http.Dir("./assets/static"))

	// register file server as handler for all URL paths starting with "/static/"
	// strip "/static" prefix before request reaches file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// middleware chain for dynamic application routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF)

	// register handler functions and route patterns
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// middleware choain for protected (authenticated-only) application routes
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// middleware chain containing 'standard' middleware
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
