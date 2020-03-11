package middlewares

import (
	"net/http"
	"turkdev/app/src/templates"
	"turkdev/shared"
)

/*Auth checks if user is authenticated to process the request.*/
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated, _, err := shared.IsAuthenticated(r)
		if err != nil {
			err = templates.RenderInLayout(w, "app/src/templates/errors/500.html", nil)
			if err != nil {
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
