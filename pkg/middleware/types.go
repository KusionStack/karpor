package middleware

// contextKey is a type used to define keys for context values. The name
// property is used to uniquely identify the context value.
type contextKey struct {
	name string // name is the identifier for the context value.
}
