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
	"context"
	"fmt"

	"github.com/pkg/errors"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/yaml"

	"github.com/KusionStack/karpor/config"
)

// ConvertToUnstructured converts the structured object to unstructured
func ConvertToUnstructured(obj runtime.Object) (*unstructured.Unstructured, error) {
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to Unstructured: %v", err)
	}

	return &unstructured.Unstructured{
		Object: unstructuredObj,
	}, nil
}

// ApplyCrds applies crds to user cluster before other resources.
func ApplyCrds(ctx context.Context, dynamicClient dynamic.Interface) error {
	for _, crd := range config.CrdList {
		var objMap map[string]interface{}
		err := yaml.Unmarshal(crd, &objMap)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{
			Object: objMap,
		}
		err = CreateOrUpdateUnstructured(ctx, dynamicClient, apiextensionsv1.SchemeGroupVersion.WithResource("customresourcedefinitions"), "", unstructuredObj)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateOrUpdateUnstructured creates or updates object using dynamic client.
func CreateOrUpdateUnstructured(ctx context.Context, dynamicClient dynamic.Interface, gvr schema.GroupVersionResource, namespace string, newObject *unstructured.Unstructured) error {
	resourceClient := dynamicClient.Resource(gvr).Namespace(namespace)

	existingObj, getErr := resourceClient.Get(ctx, newObject.GetName(), metav1.GetOptions{})
	if getErr != nil {
		if apierrors.IsNotFound(getErr) {
			// set initial resource version
			newObject.SetResourceVersion("0")
			_, createErr := resourceClient.Create(ctx, newObject, metav1.CreateOptions{})
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create resource")
			}
		} else {
			return errors.Wrapf(getErr, "failed to get resource: %v", getErr)
		}
	} else {
		// set uid and resource version for existed object
		newObject.SetResourceVersion(existingObj.GetResourceVersion())
		newObject.SetUID(existingObj.GetUID())

		_, updateErr := resourceClient.Update(ctx, newObject, metav1.UpdateOptions{})
		if updateErr != nil {
			return errors.Wrapf(updateErr, "failed to update resource: %v", newObject.GetName())
		}
	}

	return nil
}
