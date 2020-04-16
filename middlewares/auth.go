package middlewares

import (
	"net/http"
	"turkdev/shared"
	"turkdev/templates"

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
