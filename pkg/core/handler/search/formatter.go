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

package search

import (
	"fmt"
	"reflect"
	"strings"

	jp "github.com/KusionStack/karpor/pkg/util/jsonpath"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/jsonpath"
)

type Formatter interface {
	Format(obj runtime.Object) (any, error)
}

// ParseObjectFormatter parses a format string and returns a Formatter.
func ParseObjectFormatter(format string) (Formatter, error) {
	parts := strings.SplitN(format, "=", 2)
	if parts[0] == "" || parts[0] == "origin" {
		return &NopFormatter{}, nil
	}

	spec := ""
	if len(parts) > 1 {
		spec = parts[1]
	}

	switch parts[0] {
	case "custom-columns":
		return NewCustomColumnsFormatter(spec)
	default:
		return nil, fmt.Errorf("unsupported format: %s", parts[0])
	}
}

type NopFormatter struct{}

// Format keeps the obj unchanged.
func (f *NopFormatter) Format(obj runtime.Object) (any, error) {
	return obj, nil
}

// NewCustomColumnsFormatter creates a custom columns formatter from a comma separated list of <header>:<jsonpath-field-spec> pairs.
// e.g. NAME:metadata.name,API_VERSION:apiVersion creates a formatter that formats the input to:
//
// {"fields":{"API_VERSION":"v1","NAME":"foo"},"titles":["NAME","API_VERSION"]}
func NewCustomColumnsFormatter(spec string) (*customColumnsFormatter, error) {
	if len(spec) == 0 {
		return nil, fmt.Errorf("custom-columns format specified but no custom columns given")
	}

	parts := strings.Split(spec, ",")
	parsers := make([]*jsonpath.JSONPath, len(parts))
	titles := make([]string, len(parts))

	for ix := range parts {
		colSpec := strings.SplitN(parts[ix], ":", 2)
		if len(colSpec) != 2 {
			return nil, fmt.Errorf("unexpected custom-columns spec: %s, expected <header>:<json-path-expr>", parts[ix])
		}

		spec, err := jp.RelaxedJSONPathExpression(colSpec[1])
		if err != nil {
			return nil, err
		}

		parsers[ix] = jsonpath.New(fmt.Sprintf("column%d", ix)).AllowMissingKeys(true)
		if err := parsers[ix].Parse(spec); err != nil {
			return nil, err
		}

		titles[ix] = colSpec[0]
	}

	return &customColumnsFormatter{titles: titles, parsers: parsers}, nil
}

type customColumnsFormatter struct {
	titles  []string
	parsers []*jsonpath.JSONPath
}

type CustomColumnsOutput struct {
	Fields map[string]any `json:"fields,omitempty"`

	// Titles is to imply the keys order as the keys order of the fields is random.
	Titles []string `json:"titles,omitempty"`
}

// Format likes the `kubectl get -o 'custom-columns=<spec>'`, extracts and returns specified fields from input.
func (f *customColumnsFormatter) Format(obj runtime.Object) (any, error) {
	fields := map[string]any{}

	for ix, parser := range f.parsers {
		var values [][]reflect.Value
		var err error
		if unstructured, ok := obj.(runtime.Unstructured); ok {
			values, err = parser.FindResults(unstructured.UnstructuredContent())
		} else {
			values, err = parser.FindResults(reflect.ValueOf(obj).Elem().Interface())
		}

		if err != nil {
			return nil, err
		}

		if len(values) == 0 || len(values[0]) == 0 {
			fields[f.titles[ix]] = nil
			continue
		}

		typed := make([]any, len(values[0]))
		for valIx, val := range values[0] {
			typed[valIx] = val.Interface()
		}

		if len(typed) == 1 {
			fields[f.titles[ix]] = typed[0]
		} else {
			fields[f.titles[ix]] = typed
		}
	}

	return CustomColumnsOutput{Fields: fields, Titles: f.titles}, nil
}
