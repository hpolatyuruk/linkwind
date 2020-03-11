package middlewares

import (
	"net/http"
	"turkdev/app/src/templates"
)

/*Error is a middleware which handles errors for fiven http handlers.*/
func Error(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			err = templates.RenderInLayout(w, "app/src/templates/errors/500.html", nil)
			if err != nil {
				http.Error(w, "Unexpected error!", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
