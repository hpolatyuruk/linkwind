package middlewares

import (
	"linkwind/app/shared"
	"net/http"
	"strings"
)

/*NotFoundMiddleware is a middleware which handeles not found page for given handler.*/
func NotFoundMiddleware(paths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			urlPath := r.URL.Path

			if strings.HasPrefix(urlPath, shared.StaticFolderPath) ||
				pathExists(paths, urlPath) {

				next.ServeHTTP(w, r)

			} else if urlPath == "/robots.txt" {

				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("User-agent: *\nDisallow: /"))

			} else {
				shared.ReturnNotFoundTemplate(w)
				return
			}
		}
		return http.HandlerFunc(fn)
	}
}

func pathExists(paths []string, urlPath string) bool {
	for _, path := range paths {
		if urlPath == path {
			return true
		}
	}
	return false
}
