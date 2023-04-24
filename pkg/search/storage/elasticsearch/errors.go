package elasticsearch

import "fmt"

// ESError is an error type which represents a single ES error
type ESError struct {
	StatusCode int
	Message    string
}

func (e *ESError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.StatusCode, e.Message)
}
