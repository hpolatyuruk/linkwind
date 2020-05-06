package middlewares

import (
	"linkwind/app/shared"
	"linkwind/app/templates"
	"net/http"

	"github.com/getsentry/sentry-go"
)

/*Auth checks if user is authenticated to process the request.*/
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated, _, err := shared.IsAuthenticated(r)
		if err != nil {
			err = templates.RenderInLayout(w, "app/src/templates/errors/500.html", nil)
			if err != nil {
				sentry.CaptureException(err)
				http.Error(w, "Unexpected error!", http.StatusInternalServerError)
				return
			}
			return
		}
		if !isAuthenticated {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/*AuthMiddleWare checks if user is authenticated to process the request.*/
func AuthMiddleWare(authorizedPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			isPathAuthorized := false

			for _, path := range authorizedPaths {
				if r.URL.Path == path {
					isPathAuthorized = true
				}
			}

			if isPathAuthorized {
				isAuthenticated, _, err := shared.IsAuthenticated(r)
				if err != nil {
					panic(err)
				}
				if !isAuthenticated {
					http.Redirect(w, r, "/signin", http.StatusSeeOther)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
