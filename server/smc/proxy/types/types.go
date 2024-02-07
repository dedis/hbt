package types

// HTTPError defines the standard error format
type HTTPError struct {
	Title   string
	Code    uint
	Message string
	Args    map[string]interface{}
}
