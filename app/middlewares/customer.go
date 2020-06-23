package middlewares

import (
	"context"
	"fmt"
	"linkwind/app/caching"
	"linkwind/app/data"
	"linkwind/app/shared"
	"net"
	"net/http"
	"strings"
	"sync"
)

var mutex = &sync.Mutex{}

var defaultCustomerCtx = &caching.CustomerCtx{
	ID:       shared.DefaultCustomerID,
	Platform: shared.DefaultCustomerName,
	Logo:     "", // TODO: set default logo here
	Title:    "",
}

/*CustomerMiddleware sets requested customer info to request context*/
func CustomerMiddleware() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			if isLocalHost(r) {
				nexWithContext(next, w, r, defaultCustomerCtx)
				return
			}
			if isCustomDomain(r.Host) {
				handleCustomDomains(next, w, r)
				return
			}
			handleSubDomains(next, w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func handleCustomDomains(next http.Handler, w http.ResponseWriter, r *http.Request) {

	domain := r.Host

	if caching.ExistsCustomer(domain) == false {
		mutex.Lock()
		defer mutex.Unlock()
		if caching.ExistsCustomer(domain) == false {
			customer, err := data.GetCustomerByDomain(domain)
			if err != nil {
				panic(err)
			}
			if customer == nil {
				shared.ReturnNotFoundTemplate(w)
				return
			}
			var imageasB64 string
			if customer.LogoImage != nil {
				imageasB64, err = shared.EncodeLogoImageToBase64(customer.LogoImage)
				if err != nil {
					panic(err)
				}
			}
			customerCtx := &caching.CustomerCtx{
				ID:       customer.ID,
				Platform: customer.Name,
				Logo:     imageasB64,
				Title:    customer.Title,
			}
			caching.SetCustomer(domain, customerCtx)
		}
	}
	nexWithContext(next, w, r, caching.GetCustomer(domain))
}

func handleSubDomains(next http.Handler, w http.ResponseWriter, r *http.Request) {

	custName := parseCustomerName(r.Host)

	//
	// app is our default sub domain.
	// so we don't return not found if it is app
	//

	if custName == "app" {
		path := strings.ToLower(r.URL.Path)

		if isStaticPath(path) == false &&
			path != "/customer-signup" &&
			path != "/exists-custom-domain" {
			shared.ReturnNotFoundTemplate(w)
			return
		}
		nexWithContext(next, w, r, defaultCustomerCtx)
		return
	}

	if caching.ExistsCustomer(custName) == false {
		mutex.Lock()
		defer mutex.Unlock()
		if caching.ExistsCustomer(custName) == false {
			customer, err := data.GetCustomerByName(custName)
			if err != nil {
				panic(err)
			}
			if customer == nil {
				shared.ReturnNotFoundTemplate(w)
				return
			}

			var imageasB64 string
			if customer.LogoImage != nil {
				imageasB64, err = shared.EncodeLogoImageToBase64(customer.LogoImage)
				if err != nil {
					panic(err)
				}
			}
			customerCtx := &caching.CustomerCtx{
				ID:       customer.ID,
				Logo:     imageasB64,
				Platform: customer.Name,
				Title:    customer.Title,
			}
			caching.SetCustomer(custName, customerCtx)
		}
	}
	nexWithContext(next, w, r, caching.GetCustomer(custName))
}

func nexWithContext(next http.Handler, w http.ResponseWriter, r *http.Request, customersOBJ *caching.CustomerCtx) {
	ctx := context.WithValue(r.Context(), shared.CustomerContextKey, customersOBJ)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func parseCustomerName(host string) string {
	index := strings.Index(host, ".")
	if index < 0 {
		panic(fmt.Errorf("Unexpected host format %s", host))
	}
	return host[0:index]
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

func isLocalHost(r *http.Request) bool {
	ipAddress, err := getIPFromRequest(r)
	if err != nil {
		panic(err)
	}
	return ipAddress.String() == "127.0.0.1" ||
		ipAddress.String() == "::1" ||
		ipAddress.String() == "localhost"
}

func isStaticPath(path string) bool {
	return strings.Index(path, shared.StaticFolderPath) > -1
}

func isCustomDomain(host string) bool {
	return strings.Index(strings.ToLower(host), "linkwind.co") == -1
}
