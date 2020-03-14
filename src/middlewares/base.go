package middlewares

import (
	"net/http"
)

// Middleware (this function) makes adding more than one layer of middleware easy
// by specifying them as a list. It will run the last specified handler first.
// For more: https://golangcode.com/middleware-on-handlers/
func Middleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
