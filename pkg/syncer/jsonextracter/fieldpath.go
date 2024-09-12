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

package jsonextracter

import (
	"fmt"
)

// NewNestedFieldPath constructs a FieldPathExtracter.
func NewNestedFieldPath(nestedField []string, allowMissingKeys bool) *NestedFieldPath {
	return &NestedFieldPath{nestedField: nestedField, allowMissingKeys: allowMissingKeys}
}

// NestedFieldPath is used to wrap NestedFieldNoCopy function as an Extracter.
type NestedFieldPath struct {
	nestedField      []string
	allowMissingKeys bool
}

// Extract outputs the nestedField's value and its upstream structure.
func (f *NestedFieldPath) Extract(data map[string]interface{}) (map[string]interface{}, error) {
	return NestedFieldNoCopy(data, f.allowMissingKeys, f.nestedField...)
}

// NestedFieldNoCopy is similar to JSONPath.Extract. The difference is that it
// can only operate on map and does not support list, but has better performance.
func NestedFieldNoCopy(data map[string]interface{}, allowMissingKeys bool, fields ...string) (map[string]interface{}, error) {
	if len(fields) == 0 {
		return nil, nil
	}

	result := map[string]interface{}{}
	cur := result

	for i, field := range fields {
		if val, ok := data[field]; ok {
			if i != len(fields)-1 {
				if data, ok = val.(map[string]interface{}); !ok {
					return nil, fmt.Errorf("%v is of the type %T, expected map[string]interface{}", val, val)
				}

				m := map[string]interface{}{}
				cur[field] = m
				cur = m
			} else {
				cur[field] = val
			}
		} else {
			if allowMissingKeys {
				return result, nil
			}
			return nil, fmt.Errorf("field %q not exist", field)
		}
	}

	return result, nil
}
