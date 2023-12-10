package audit

// Payload represents the structure for audit request data. It includes the
// manifest which is typically a string containing declarative configuration
// data that needs to be audited.
type Payload struct {
	Manifest string `json:"manifest"` // Manifest is the content to be audited.
}
