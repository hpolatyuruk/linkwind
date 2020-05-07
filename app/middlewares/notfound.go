package middlewares

import (
	"linkwind/app/shared"
	"net/http"

	"github.com/getsentry/sentry-go"
)

/*NotFoundMiddleWare is a middleware which handeles not found page for given handler.*/
func NotFoundMiddleWare(paths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			urlPath := r.URL.Path

			if urlPath == "/robots.txt" {

				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("User-agent: *\nDisallow: /"))

			} else if pathExists(paths, urlPath) == false {

				renderFile(w, "templates/errors/404.html", http.StatusNotFound)

			} else {
				// Call the next handler
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}

func pathExists(paths []string, urlPath string) bool {
	for _, path := range paths {
		if path == urlPath {
			return true
		}
	}
	return false
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
