package middlewares

import (
	"linkwind/app/shared"
	"net/http"
)

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
