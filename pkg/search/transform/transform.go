// Copyright The Karbour Authors.
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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var defaultRegistry = NewRegistry()

func init() {
	Register("patch", Patch)
	Register("replace", Replace)
}

type TransformFunc func(original interface{}, transformText string) (target interface{}, err error)

func Register(tType string, transFunc TransformFunc) {
	defaultRegistry.Register(tType, transFunc)
}

func GetTransformFunc(transformerType string) (TransformFunc, bool) {
	return defaultRegistry.Get(transformerType)
}

type TransformFuncRegistry struct {
	transformers map[string]TransformFunc
}

func NewRegistry() *TransformFuncRegistry {
	return &TransformFuncRegistry{transformers: make(map[string]TransformFunc)}
}

func (r *TransformFuncRegistry) Register(tType string, transFunc TransformFunc) {
	r.transformers[tType] = transFunc
}

func (r *TransformFuncRegistry) Get(transformerType string) (transFunc TransformFunc, found bool) {
	transFunc, found = r.transformers[transformerType]
	return
}

func Patch(original interface{}, patchText string) (interface{}, error) {
	patch, err := jsonpatch.DecodePatch([]byte(patchText))
	if err != nil {
		return nil, fmt.Errorf("patch is invalid: %v", err)
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
		return nil, fmt.Errorf(`type "%vT not supported`, original)
	}
	dest = u.NewEmptyInstance()

	if err := json.Unmarshal(modifiedJSON, &dest); err != nil {
		return nil, fmt.Errorf("json decoding error: %v", err)
	}
	return &dest, nil
}

func Replace(original interface{}, jsonString string) (interface{}, error) {
	var dest unstructured.Unstructured
	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		return nil, fmt.Errorf("json decoding error: %v", err)
	}
	return &dest, nil
}
