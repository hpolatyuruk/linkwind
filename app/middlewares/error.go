package middlewares

import (
	"errors"
	"fmt"
	"linkwind/app/shared"
	"net/http"

	"github.com/getsentry/sentry-go"
)

/*Error is a middleware which handles errors for fiven http handlers.*/
func Error(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			sentry.CaptureException(err)
			byteValue, err := shared.ReadFile("templates/errors/500.html")
			if err != nil {
				http.Error(w, "Unexpected error!", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(byteValue)
		}
	})
}

/*ErrorMiddleWare is a middleware which handles errors for fiven http handlers.*/
func ErrorMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					fmt.Println("Recovered in f", recovered)
					// find out exactly what the error was and set err
					switch recoveredType := recovered.(type) {
					case string:
						recovered = errors.New(recoveredType)
					case error:
						recovered = recoveredType
					default:
						// Fallback err (per specs, error strings should be lowercase w/o punctuation
						recovered = errors.New("unknown panic")
					}
					sentry.CaptureException(recovered.(error))
					w.WriteHeader(http.StatusInternalServerError)
					byteValue, err := shared.ReadFile("templates/errors/500.html")
					if err != nil {
						http.Error(w, "Unexpected error!", http.StatusInternalServerError)
					}
					w.Write(byteValue)
				}
			}()

			// Call the next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
