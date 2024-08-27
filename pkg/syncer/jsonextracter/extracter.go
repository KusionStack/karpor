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

	"k8s.io/client-go/util/jsonpath"
)

type Extracter interface {
	Extract(data map[string]interface{}) (map[string]interface{}, error)
}

// BuildExtracter automatically determines whether to use FieldPathExtracter or JSONPathExtracter.
// If the input jsonPath only involves map operations, it will return FieldPathExtracter,
// as it has better performance.
func BuildExtracter(jsonPath string, allowMissingKeys bool) (Extracter, error) {
	parser, err := Parse(jsonPath, jsonPath)
	if err != nil {
		return nil, fmt.Errorf("error in parsing path %q: %w", jsonPath, err)
	}

	rootNodes := parser.Root.Nodes
	if len(rootNodes) == 0 {
		return NewNestedFieldPath(nil, allowMissingKeys), nil
	}

	if len(rootNodes) == 1 {
		nodes := rootNodes[0].(*jsonpath.ListNode).Nodes
		fields := make([]string, 0, len(nodes))
		for _, node := range nodes {
			if node.Type() == jsonpath.NodeField {
				fields = append(fields, node.(*jsonpath.FieldNode).Value)
			}
		}

		if len(nodes) == len(fields) {
			return NewNestedFieldPath(fields, allowMissingKeys), nil
		}
	}

	jp := &JSONPath{name: parser.Name, parser: parser, allowMissingKeys: allowMissingKeys}
	return jp, nil
}
