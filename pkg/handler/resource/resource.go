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

package resource

import (
	"encoding/json"
	"net/http"

	"github.com/KusionStack/karbour/pkg/apis/search"
	"github.com/KusionStack/karbour/pkg/controller/resource"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/registry"
	searchstorage "github.com/KusionStack/karbour/pkg/registry/search"

	"github.com/go-chi/chi/v5"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/server"
)

func Get(resourceCtrl *resource.ResourceController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := BuildResourceFromParam(r)
		dynamicClient, _ := multicluster.BuildDynamicClient(r.Context(), c.LoopbackClientConfig, res.Cluster)
		resourceUnstructured, _ := resourceCtrl.GetResource(r.Context(), dynamicClient, res)
		result, _ := json.MarshalIndent(resourceUnstructured, "", "  ")
		w.Write(result)
	}
}

func GetYAML(resourceCtrl *resource.ResourceController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := BuildResourceFromParam(r)
		client, _ := multicluster.BuildDynamicClient(r.Context(), c.LoopbackClientConfig, res.Cluster)
		result, _ := resourceCtrl.GetYAMLForResource(r.Context(), client, res)
		w.Write(result)
	}
}

func GetTopology(resourceCtrl *resource.ResourceController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := BuildResourceFromParam(r)
		dynamicClient, _ := multicluster.BuildDynamicClient(r.Context(), c.LoopbackClientConfig, res.Cluster)
		discoveryClient, _ := multicluster.BuildDiscoveryClient(r.Context(), c.LoopbackClientConfig, res.Cluster)
		topologyMap, _ := resourceCtrl.GetTopologyForResource(r.Context(), dynamicClient, discoveryClient, res)
		result, _ := json.MarshalIndent(topologyMap, "", "  ")
		w.Write(result)
	}
}

func SearchForResource(resourceCtrl *resource.ResourceController, c *registry.ExtraConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchQuery := r.URL.Query().Get("query")
		searchPattern := r.URL.Query().Get("pattern")
		storage := searchstorage.RESTStorageProvider{
			SearchStorageType:      c.SearchStorageType,
			ElasticSearchAddresses: c.ElasticSearchAddresses,
			ElasticSearchName:      c.ElasticSearchName,
			ElasticSearchPassword:  c.ElasticSearchPassword,
		}
		searchStorageGetter, err := storage.SearchStorageGetter()
		if err != nil {
			return
		}
		searchStorage, err := searchStorageGetter.GetSearchStorage()
		if err != nil {
			return
		}
		res, err := searchStorage.Search(r.Context(), searchQuery, searchPattern)
		if err != nil {
			return
		}

		rt := &search.UniResourceList{}
		for _, resource := range res.Resources {
			unObj := &unstructured.Unstructured{}
			unObj.SetUnstructuredContent(resource.Object)
			rt.Items = append(rt.Items, unObj)
		}
		result, _ := json.MarshalIndent(rt.Items, "", "  ")
		w.Write(result)
	}
}

func BuildResourceFromParam(r *http.Request) *resource.Resource {
	res := resource.Resource{
		Cluster:    chi.URLParam(r, "clusterName"),
		APIVersion: chi.URLParam(r, "apiVersion"),
		Kind:       chi.URLParam(r, "kind"),
		Namespace:  chi.URLParam(r, "namespaceName"),
		Name:       chi.URLParam(r, "resourceName"),
	}
	return &res
}
