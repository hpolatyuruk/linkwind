package shared

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key string

const (
	// DefaultCustomerID represents the default identifier of customer
	DefaultCustomerID int = 1
	// CustomerIDContextKey represents the key to get customer id from request context
	CustomerIDContextKey key = "CustomerID" // default

	// StaticFolderPath represents the folder that contains static files
	StaticFolderPath = "/public/"
)
