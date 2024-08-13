// Copyright The Karpor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"sync"

	"k8s.io/client-go/util/jsonpath"
)

// JSONPathParser is a struct that holds a cache of parsed JSONPath expressions for efficient reuse.
type JSONPathParser struct {
	sync.Mutex
	cache map[string]*jsonpath.JSONPath
}

var DefaultJSONPathParser = NewJSONPathParser()

// NewJSONPathParser creates and returns a new instance of JSONPathParser with an empty cache.
func NewJSONPathParser() *JSONPathParser {
	return &JSONPathParser{
		cache: make(map[string]*jsonpath.JSONPath),
	}
}

// Parse takes a JSONPath expression string and returns a parsed JSONPath object or an error if the expression is invalid.
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

// JSONPathFields is a struct that holds a JSONPathParser instance and the data to be queried by JSONPath expressions.
type JSONPathFields struct {
	jpParser *JSONPathParser
	data     interface{}
}

// NewJSONPathFields creates and returns a new instance of JSONPathFields with the given parser and data.
func NewJSONPathFields(jpParser *JSONPathParser, data interface{}) *JSONPathFields {
	return &JSONPathFields{
		jpParser: jpParser,
		data:     data,
	}
}

// Has checks if the given JSONPath expression exists in the data, returning a boolean indicating its existence.
func (fs JSONPathFields) Has(fieldPath string) (exists bool) {
	jp, err := fs.jpParser.Parse(fieldPath)
	if err != nil {
		return false
	}
	vals, err := jp.FindResults(fs.data)
	if err != nil || len(vals) == 0 || len(vals[0]) == 0 {
		return false
	}
	return len(vals) > 0
}

// Get retrieves the value at the specified JSONPath expression from the data and returns it as a string.
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
