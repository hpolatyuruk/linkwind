package middlewares

import (
	"context"
	"fmt"
	"linkwind/app/data"
	"linkwind/app/shared"
	"net"
	"net/http"
	"strings"
	"sync"
)

var customers map[string]int = map[string]int{}

/*CustomerMiddleware sets requested customer info to request context*/
func CustomerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var exists bool = false
			var mutex = &sync.Mutex{}
			var customerName string
			var customerID int = shared.DefaultCustomerID

			if isLocalHost(r) {
				nexWithContext(next, w, r, customerID)
				return
			}
			customerName = parseCustomerName(r.Host)

			//
			// app is our default domain.
			// so we don't return not fund if it is app
			//

			if customerName == "app" {

				path := strings.ToLower(r.URL.Path)

				if isStaticPath(path) == false && path != "/customer-signup" {
					shared.ReturnNotFoundTemplate(w)
					return
				}
				nexWithContext(next, w, r, customerID)
				return
			}

			exists, customerID = existsInCache(customerName)
			if exists == false {
				mutex.Lock()
				exists, customerID = existsInCache(customerName)
				if exists == false {
					customer, err := data.GetCustomerByName(customerName)
					if err != nil {
						panic(err)
					}
					if customer == nil {
						shared.ReturnNotFoundTemplate(w)
						return
					}
					customerID = customer.ID
					customers[customer.Name] = customerID
				}
				mutex.Unlock()
			}
			nexWithContext(next, w, r, customerID)
		}
		return http.HandlerFunc(fn)
	}
}

func nexWithContext(next http.Handler, w http.ResponseWriter, r *http.Request, customerID int) {
	ctx := context.WithValue(r.Context(), shared.CustomerIDContextKey, customerID)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func existsInCache(name string) (bool, int) {
	if id, ok := customers[name]; ok {
		return ok, id
	}
	return false, 0
}

func parseCustomerName(host string) string {
	index := strings.Index(host, ".")
	if index < 0 {
		panic(fmt.Errorf("Unexpected host format %s", host))
	}
	return host[0:index]
}

func isLocalHost(r *http.Request) bool {
	ipAddress, err := getIPFromRequest(r)
	if err != nil {
		panic(err)
	}
	return ipAddress.String() == "127.0.0.1" ||
		ipAddress.String() == "::1" ||
		ipAddress.String() == "localhost"
}

func getIPFromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

func isStaticPath(path string) bool {
	return strings.Index(path, shared.StaticFolderPath) > -1
}
