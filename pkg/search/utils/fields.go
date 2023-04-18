package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"sync"

	"k8s.io/client-go/util/jsonpath"
)

type JSONPathParser struct {
	sync.Mutex
	cache map[string]*jsonpath.JSONPath
}

func NewJSONPathParser() *JSONPathParser {
	return &JSONPathParser{
		cache: make(map[string]*jsonpath.JSONPath),
	}
}

func (j *JSONPathParser) Parse(path string) (*jsonpath.JSONPath, error) {
	j.Lock()
	defer j.Unlock()

	if cache, found := j.cache[path]; found {
		return cache, nil
	}

	p := jsonpath.New("fieldpath: " + path).AllowMissingKeys(true)
	if err := p.Parse(path); err != nil {
		return nil, err
	}
	j.cache[path] = p
	return p, nil
}

type JSONPathFields struct {
	jpParser *JSONPathParser
	data     interface{}
}

func NewJSONPathFields(jpParser *JSONPathParser, data interface{}) *JSONPathFields {
	return &JSONPathFields{
		jpParser: jpParser,
		data:     data,
	}
}

func (fs JSONPathFields) Has(fieldPath string) (exists bool) {
	jp, err := fs.jpParser.Parse(fieldPath)
	if err != nil {
		return false
	}
	vals, err := jp.FindResults(fs.data)
	if err != nil {
		return false
	}
	return len(vals) > 0
}

func (fs JSONPathFields) Get(fieldPath string) (value string) {
	fieldPath, err := RelaxedJSONPathExpression(fieldPath)
	if err != nil {
		return ""
	}

	jp, err := fs.jpParser.Parse(fieldPath)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	if err := jp.Execute(&buf, fs.data); err != nil {
		return ""
	}
	return buf.String()
}

var jsonRegexp = regexp.MustCompile(`^\{\.?([^{}]+)\}$|^\.?([^{}]+)$`)

// copied from kubectl
func RelaxedJSONPathExpression(pathExpression string) (string, error) {
	if len(pathExpression) == 0 {
		return pathExpression, nil
	}
	submatches := jsonRegexp.FindStringSubmatch(pathExpression)
	if submatches == nil {
		return "", fmt.Errorf("unexpected path string, expected a 'name1.name2' or '.name1.name2' or '{name1.name2}' or '{.name1.name2}'")
	}
	if len(submatches) != 3 {
		return "", fmt.Errorf("unexpected submatch list: %v", submatches)
	}
	var fieldSpec string
	if len(submatches[1]) != 0 {
		fieldSpec = submatches[1]
	} else {
		fieldSpec = submatches[2]
	}
	return fmt.Sprintf("{.%s}", fieldSpec), nil
}
