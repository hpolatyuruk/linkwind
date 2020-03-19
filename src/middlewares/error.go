package middlewares

import (
	"net/http"
	"turkdev/src/shared"

	"github.com/getsentry/sentry-go"
)

/*Error is a middleware which handles errors for fiven http handlers.*/
func Error(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			sentry.CaptureException(err)
			byteValue, err := shared.ReadFile("src/templates/errors/500.html")
			if err != nil {
				http.Error(w, "Unexpected error!", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(byteValue)
		}
	})
}
