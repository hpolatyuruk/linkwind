package shared

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key string

const (
	// DefaultCustomerID represents the default identifier of customer
	DefaultCustomerID int = 1
	// DefaultCustomerName represents the default name of customer
	DefaultCustomerName string = "demo"
	// CustomerContextKey represents the key to get customer id from request context
	CustomerContextKey key = "CustomerID" // default

	// UserContextKey represents the key to get authenticated user from request context
	UserContextKey key = "UserID" // default

	// StaticFolderPath represents the folder that contains static files
	StaticFolderPath = "/public/"
)
