package middlewares

import (
	"linkwind/app/shared"
	"net/http"

	"github.com/getsentry/sentry-go"
)

/*NotFound is a middleware which handeles not found page for given handler.*/
func NotFound(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			err := f(w, r)
			if err != nil {
				sentry.CaptureException(err)
				renderFile(w, "src/templates/errors/500.html", http.StatusInternalServerError)
			}
		} else if r.URL.Path == "/robots.txt" {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("User-agent: *\nDisallow: /"))
		} else {
			renderFile(w, "src/templates/errors/404.html", http.StatusNotFound)
		}
	}
}
func renderFile(w http.ResponseWriter, path string, statusCode int) {
	byteValue, err := shared.ReadFile(path)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, "Unexpected error!", http.StatusInternalServerError)
	}
	w.WriteHeader(statusCode)
	w.Write(byteValue)
}
