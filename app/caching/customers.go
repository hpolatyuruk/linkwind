package caching

/*CustomerCtx respresents the customer object in the request context*/
type CustomerCtx struct {
	ID       int
	Platform string
	Logo     string
	Title    string
}

var customers map[string]*CustomerCtx = map[string]*CustomerCtx{}

/*ExistsCustomer checks whether customer context object exists in cache or not*/
func ExistsCustomer(custNameOrDomain string) bool {
	if _, ok := customers[custNameOrDomain]; ok {
		return ok
	}
	return false
}

/*GetCustomer get customer context object from cache*/
func GetCustomer(custNameOrDomain string) *CustomerCtx {
	customer, _ := customers[custNameOrDomain]
	return customer
}

/*SetCustomer sets customer context object in cache*/
func SetCustomer(custNameOrDomain string, customer *CustomerCtx) {
	customers[custNameOrDomain] = customer
}

/*DeleteCustomer deletes customer context object from cache*/
func DeleteCustomer(custNameOrDomain string) {
	delete(customers, custNameOrDomain)
}
