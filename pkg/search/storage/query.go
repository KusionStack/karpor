package storage

type Operator string

const (
	Equals Operator = "="
)

type Query struct {
	Key      string
	Values   []string
	Operator Operator
}
