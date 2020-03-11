package middlewares

import (
	"net/http"
	"turkdev/app/src/templates"
)

/*NotFound is a middleware which handeles not found page for given handler.*/
func NotFound(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			err := f(w, r)
			if err != nil {
				renderInternalServerError(w)
			}
		} else if r.URL.Path == "/robots.txt" {
			w.Header().Set("Content-Type", "text/plain")
			// TODO: Should be removed in production.
			w.Write([]byte("User-agent: *\nDisallow: /"))
		} else {
			renderNotFound(w)
		}
	}
}
func renderInternalServerError(w http.ResponseWriter) {
	err := templates.RenderInLayout(w, "app/src/templates/errors/500.html", nil)
	if err != nil {
		http.Error(w, "Unexpected error!", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func renderNotFound(w http.ResponseWriter) {
	err := templates.RenderInLayout(w, "app/src/templates/errors/404.html", nil)
	if err != nil {
		http.Error(w, "Unexpected error!", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNotFound)
}
