/*
Copyright The Karbour Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"context"
	"fmt"

	"github.com/KusionStack/karbour/pkg/apis/cluster"
	clusterv1beta1 "github.com/KusionStack/karbour/pkg/apis/cluster/v1beta1"
	yaml "gopkg.in/yaml.v3"

	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/klog/v2"
)

var _ rest.Lister = &YAMLRest{}

type YAMLRest struct {
	Store *genericregistry.Store
}

// New returns empty Cluster object.
func (r *YAMLRest) New() runtime.Object {
	return &cluster.Cluster{}
}

func (r *YAMLRest) NewList() runtime.Object {
	return &cluster.ClusterList{}
}

// Destroy cleans up resources on shutdown.
func (r *YAMLRest) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *YAMLRest) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.Store.Get(ctx, name, options)
}

func (r *YAMLRest) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	rt := &cluster.Cluster{}
	info, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve request info from context")
	}
	klog.Infof("Getting YAML for cluster %s...", info.Name)

	// Locate the cluster resource based on the name in the request
	obj, err := r.Store.Get(ctx, info.Name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	obj.GetObjectKind().SetGroupVersionKind(clusterv1beta1.SchemeGroupVersion.WithKind("Cluster"))
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	objYAML, err := yaml.Marshal(unstructuredObj)
	if err != nil {
		return nil, err
	}
	rt.Status.YAMLString = string(objYAML)
	return rt, nil
}

func (r *YAMLRest) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.Store.ConvertToTable(ctx, object, tableOptions)
}
