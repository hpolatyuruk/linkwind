package middlewares

import (
	"linkwind/app/shared"
	"net/http"

	"github.com/getsentry/sentry-go"
)

/*NotFoundMiddleWare is a middleware which handeles not found page for given handler.*/
func NotFoundMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" || r.URL.Path == "/index.html" {
				// Call the next handler
				next.ServeHTTP(w, r)
			} else if r.URL.Path == "/robots.txt" {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("User-agent: *\nDisallow: /"))
			} else {
				renderFile(w, "templates/errors/404.html", http.StatusNotFound)
			}
		}
		return http.HandlerFunc(fn)
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
