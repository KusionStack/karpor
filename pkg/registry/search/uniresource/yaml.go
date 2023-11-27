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

package uniresource

import (
	"context"
	"fmt"

	cluster "github.com/KusionStack/karbour/pkg/apis/cluster"
	"github.com/KusionStack/karbour/pkg/apis/search"
	clusterstorage "github.com/KusionStack/karbour/pkg/registry/cluster"
	filtersutil "github.com/KusionStack/karbour/pkg/util/filters"
	topologyutil "github.com/KusionStack/karbour/pkg/util/topology"
	yaml "gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
)

type YAMLREST struct {
	Storage storage.SearchStorage
	Cluster *clusterstorage.Storage
}

// BuildDynamicClient returns a dynamic client based on the cluster name in the request
func (r *YAMLREST) BuildDynamicClient(ctx context.Context) (*dynamic.DynamicClient, error) {
	// Extract the cluster name from context
	resourceDetail, ok := filtersutil.ResourceDetailFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("name, namespace, cluster, apiVersion and kind are used to locate a unique resource so they can't be empty")
	}

	// Locate the cluster resource and build config with it
	obj, err := r.Cluster.Status.Store.Get(ctx, resourceDetail.Cluster, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clusterFromContext := obj.(*cluster.Cluster)
	klog.Infof("Cluster found: %s", clusterFromContext.Name)
	config, err := clusterstorage.NewConfigFromCluster(clusterFromContext)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster client config %s", clusterFromContext.Name)
	}

	// Create the dynamic client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// New returns an empty cluster proxy subresource.
func (r *YAMLREST) New() runtime.Object {
	return &search.UniresourceYAML{}
}

func (r *YAMLREST) NewList() runtime.Object {
	return &search.UniresourceYAMLList{}
}

// Destroy cleans up resources on shutdown.
func (r *YAMLREST) Destroy() {
	// Given that underlying store is shared with REST,
	// we don't destroy it here explicitly.
}

func (r *YAMLREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: add real logic of convert to table when the storage layer is implemented
	return rest.NewDefaultTableConvertor(search.Resource("uniresources")).ConvertToTable(ctx, object, tableOptions)
}

func (r *YAMLREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	rt := &search.UniResourceList{}
	client, err := r.BuildDynamicClient(ctx)
	if err != nil {
		return rt, err
	}
	resource, ok := filtersutil.ResourceDetailFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("name, namespace, cluster, apiVersion and kind are used to locate a unique resource so they can't be empty")
	}
	queryString := fmt.Sprintf("%s where name = '%s' AND namespace = '%s' AND cluster = '%s' AND apiVersion = '%s' AND kind = '%s'", SQLQueryDefault, resource.Name, resource.Namespace, resource.Cluster, resource.APIVersion, resource.Kind)
	// TODO: Should we enforce all fields to be present? Or do we allow topology graph for multiple (fuzzy search) resources at a time?
	// if resource.Namespace != "" {
	// 	queryString += fmt.Sprintf(" AND namespace = '%s'", resource.Namespace)
	// }
	// ...
	klog.Infof("Query string: %s", queryString)
	gvr, err := topologyutil.GetGVRFromGVK(resource.APIVersion, resource.Kind)
	if err != nil {
		return nil, err
	}
	res, err := client.Resource(gvr).Namespace(resource.Namespace).Get(ctx, resource.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	objYAML, err := yaml.Marshal(res.Object)
	if err != nil {
		panic(err.Error())
	}
	klog.Infof("---\n%s\n", string(objYAML))
	rt.Items = append(rt.Items, res)
	return rt, nil
}
