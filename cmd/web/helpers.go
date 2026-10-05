package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) decodePostForm(r *http.Request, dst any) error {
	// parse request
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// pass target destination for decoding form data
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// check target destination
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); ok {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack()) // convert byte slice to string
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends specific status code and corresponding description to user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) error {
	// retrieve template set from cache
	ts, ok := app.templateCache[page]
	if !ok {
		return fmt.Errorf("the template %s does not exist in the cache", page)
	}

	// initialize new buffer
	buf := new(bytes.Buffer)

	// write template to buffer
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)

	// write contents of buffer to http.ResponseWriter,
	_, err = buf.WriteTo(w)
	return err
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
