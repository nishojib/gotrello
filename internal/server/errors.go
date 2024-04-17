package server

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// serverError writes a log entry at Error level (including the request method and URI as
// attributes), then sends a generic 500 Internal Server Error response to the user
func serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	slog.Error(
		err.Error(),
		"method", method,
		"uri", uri,
		"trace", trace,
	)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific status code and corresponding description to the user.
func clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound sends 404 status code to the user
func notFound(w http.ResponseWriter) {
	clientError(w, http.StatusNotFound)
}

// methodNotAllowed sends 405 status code to the user
func methodNotAllowed(w http.ResponseWriter) {
	clientError(w, http.StatusMethodNotAllowed)
}
