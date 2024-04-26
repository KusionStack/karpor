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

package topology

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oliveagle/jsonpath"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
)

type Scope struct {
	Scope meta.RESTScopeName
}

func (s Scope) Name() meta.RESTScopeName {
	return s.Scope
}

// IsMapSubset returns true if sub is a subnet of m
func IsMapSubset[K, V comparable](m, sub map[K]V) bool {
	if len(sub) > len(m) {
		return false
	}
	for k, vsub := range sub {
		if vm, found := m[k]; !found || vm != vsub {
			return false
		}
	}
	return true
}

// OwnerReferencesMatch returns true if parent is in the child's ownerReferences field
func OwnerReferencesMatch(parent, child unstructured.Unstructured) (bool, error) {
	ownerReferences := child.GetOwnerReferences()
	for _, owner := range ownerReferences {
		if owner.UID == parent.GetUID() {
			return true, nil
		}
	}
	return false, nil
}

// LabelSelectorsMatch returns true if the labels set in selectorPath in the selectorObj can select selectedObj
func LabelSelectorsMatch(selectorObj, selectedObj unstructured.Unstructured, selectorPath string) (bool, error) {
	selectorSplit := strings.Split(selectorPath, ".")
	selectors, _, _ := unstructured.NestedStringMap(selectorObj.Object, selectorSplit...)
	if len(selectors) == 0 {
		return false, fmt.Errorf("shouldn't have empty selector and a selector type relationship")
	}
	labels := selectedObj.GetLabels()
	return IsMapSubset(labels, selectors), nil
}

// JSONPathMatch returns true if source.criteriaKey(name/namespace) matches target.criteriaValue(JSONPath)
// criteriaSet contains a list of map from criteriaKey to criteriaValue
// Returns true if either map returns a full match based on length of the map
func JSONPathMatch(source, target unstructured.Unstructured, criteriaSet []map[string]string) (bool, error) {
	criteriaMatchCount := 0
	for _, criteriaMap := range criteriaSet {
		for criteriaKey, criteriaValue := range criteriaMap {
			targetValue, _ := GetNestedValue(target, criteriaValue)
			var sourceValue string
			if criteriaKey == "name" {
				// match name
				sourceValue = source.GetName()
			} else if criteriaKey == "namespace" {
				// match namespace
				sourceValue = source.GetNamespace()
			} else {
				// shouldn't be anything else
				return false, fmt.Errorf("shouldn't have anything other than name or namespace")
			}
			// If targetValue is an array
			if reflect.TypeOf(targetValue).Kind() == reflect.Slice {
				// If any of the elements in the targetValue slice matches sourceValue, consider it a match
				// Example: If the secret name appears in the list of $.spec.volumes[:].secret.secretName
				for i := 0; i < reflect.ValueOf(targetValue).Len(); i++ {
					if reflect.ValueOf(targetValue).Index(i).Interface().(string) == sourceValue {
						criteriaMatchCount++
					}
				}
			} else if reflect.TypeOf(targetValue).Kind() == reflect.String {
				if sourceValue == targetValue {
					criteriaMatchCount++
				}
			}
		}
		// Only returns match if all of the matching criteria in the map are true
		if criteriaMatchCount == len(criteriaMap) {
			klog.Infof("Found a match based on JSONPath!")
			return true, nil
		}
	}
	return false, nil
}

// GetNestedValue returns nested value from the JSONPath obj.criteria
func GetNestedValue(obj unstructured.Unstructured, criteria string) (interface{}, error) {
	pat := jsonpath.MustCompile(criteria)
	res, err := pat.Lookup(obj.Object)
	if err != nil {
		return "", err
	}
	return res, nil
}

// GetGVRFromGVK retrieves the GroupVersionResource for a given API version and kind.
func GetGVRFromGVK(apiVersion, kind string) (schema.GroupVersionResource, error) {
	gv, _ := schema.ParseGroupVersion(apiVersion)
	gvk := gv.WithKind(kind)
	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{})
	scope := Scope{"namespace"}
	mapper.Add(gvk, scope)
	mapping, err := mapper.RESTMapping(schema.GroupKind{Group: gv.Group, Kind: kind}, gv.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	gvr := mapping.Resource
	return gvr, nil
}
