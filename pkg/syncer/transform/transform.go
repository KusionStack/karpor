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

package transform

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var defaultRegistry = NewRegistry()

func init() {
	Register("patch", Patch)
	Register("replace", Replace)
}

// TransformFunc is a type definition for transformation functions that take an original data structure and a transformation string, and return a transformed data structure and an error.
type TransformFunc func(original interface{}, transformText string) (target interface{}, err error)

// Register function registers a transformation function with the given type.
func Register(tType string, transFunc TransformFunc) {
	defaultRegistry.Register(tType, transFunc)
}

// GetTransformFunc retrieves a registered transformation function by type.
func GetTransformFunc(transformerType string) (TransformFunc, bool) {
	return defaultRegistry.Get(transformerType)
}

// TransformFuncRegistry is a struct that holds a map of transformation functions.
type TransformFuncRegistry struct {
	transformers map[string]TransformFunc
}

// NewRegistry creates and returns a new instance of TransformFuncRegistry.
func NewRegistry() *TransformFuncRegistry {
	return &TransformFuncRegistry{transformers: make(map[string]TransformFunc)}
}

// Register method of TransformFuncRegistry registers a transformation function with the given type.
func (r *TransformFuncRegistry) Register(tType string, transFunc TransformFunc) {
	r.transformers[tType] = transFunc
}

// Get method of TransformFuncRegistry retrieves a registered transformation function by type.
func (r *TransformFuncRegistry) Get(transformerType string) (transFunc TransformFunc, found bool) {
	transFunc, found = r.transformers[transformerType]
	return
}

// Patch function applies a JSON patch to the original data structure.
func Patch(original interface{}, patchText string) (interface{}, error) {
	patch, err := jsonpatch.DecodePatch([]byte(patchText))
	if err != nil {
		return nil, errors.Wrap(err, "patch is invalid")
	}

	originalJSON, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	modifiedJSON, err := patch.Apply(originalJSON)
	if err != nil {
		return nil, err
	}

	var dest interface{}
	u, ok := original.(runtime.Unstructured)
	if !ok {
		return nil, fmt.Errorf(`type %T not supported`, original)
	}
	dest = u.NewEmptyInstance()

	if err := json.Unmarshal(modifiedJSON, &dest); err != nil {
		return nil, errors.Wrap(err, "json decoding error")
	}
	return dest, nil
}

// Replace function replaces the original data structure with the new one derived from the JSON string.
func Replace(original interface{}, jsonString string) (interface{}, error) {
	var dest unstructured.Unstructured
	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		return nil, errors.Wrap(err, "json decoding error")
	}
	return &dest, nil
}
