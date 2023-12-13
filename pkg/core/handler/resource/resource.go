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
	"github.com/KusionStack/karbour/pkg/core/manager/resource"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/registry"
	searchstorage "github.com/KusionStack/karbour/pkg/registry/search"

	"github.com/go-chi/chi/v5"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/server"
)

// Get returns an HTTP handler function that returns a Kubernetes
// resource. It utilizes a ResourceManager to execute the logic.
//
//	@Summary		Get returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Description	This endpoint returns a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Tags			resource
//	@Produce		json
//	@Success		200	{object}	unstructured.Unstructured	"Unstructured object"
//	@Failure		400	{string}	string						"Bad Request"
//	@Failure		401	{string}	string						"Unauthorized"
//	@Failure		404	{string}	string						"Not Found"
//	@Failure		405	{string}	string						"Method Not Allowed"
//	@Failure		429	{string}	string						"Too Many Requests"
//	@Failure		500	{string}	string						"Internal Server Error"
//	@Router			/api/v1/resource/cluster/{clusterName}/{apiVersion}/namespace/{namespaceName}/{kind}/name/{resourceName}/ [get]
func Get(resourceMgr *resource.ResourceManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := BuildResourceFromParam(r)
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, res.Cluster)
		resourceUnstructured, _ := resourceMgr.GetResource(r.Context(), client, res)
		result, _ := json.MarshalIndent(resourceUnstructured, "", "  ")
		w.Write(result)
	}
}

// GetYAML returns an HTTP handler function that returns a Kubernetes
// resource YAML. It utilizes a ResourceManager to execute the logic.
//
//	@Summary		GetYAML returns a Kubernetes resource YAML by name, namespace, cluster, apiVersion and kind.
//	@Description	This endpoint returns a Kubernetes resource YAML by name, namespace, cluster, apiVersion and kind.
//	@Tags			resource
//	@Produce		json
//	@Success		200	{array}		byte	"Byte array"
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		404	{string}	string	"Not Found"
//	@Failure		405	{string}	string	"Method Not Allowed"
//	@Failure		429	{string}	string	"Too Many Requests"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/api/v1/resource/cluster/{clusterName}/{apiVersion}/namespace/{namespaceName}/{kind}/name/{resourceName}/yaml [get]
func GetYAML(resourceMgr *resource.ResourceManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := BuildResourceFromParam(r)
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, res.Cluster)
		result, _ := resourceMgr.GetYAMLForResource(r.Context(), client, res)
		w.Write(result)
	}
}

// GetTopology returns an HTTP handler function that returns a topology map for
// a Kubernetes resource. It utilizes a ResourceManager to execute the logic.
//
//	@Summary		GetTopology returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Description	This endpoint returns a topology map for a Kubernetes resource by name, namespace, cluster, apiVersion and kind.
//	@Tags			resource
//	@Produce		json
//	@Success		200	{object}	map[string]resource.ResourceTopology	"map from string to resource.ResourceTopology"
//	@Failure		400	{string}	string									"Bad Request"
//	@Failure		401	{string}	string									"Unauthorized"
//	@Failure		404	{string}	string									"Not Found"
//	@Failure		405	{string}	string									"Method Not Allowed"
//	@Failure		429	{string}	string									"Too Many Requests"
//	@Failure		500	{string}	string									"Internal Server Error"
//	@Router			/api/v1/resource/cluster/{clusterName}/{apiVersion}/namespace/{namespaceName}/{kind}/name/{resourceName}/topology [get]
func GetTopology(resourceMgr *resource.ResourceManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := BuildResourceFromParam(r)
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, res.Cluster)
		topologyMap, _ := resourceMgr.GetTopologyForResource(r.Context(), client, res)
		result, _ := json.MarshalIndent(topologyMap, "", "  ")
		w.Write(result)
	}
}

// SearchForResource returns an HTTP handler function that returns an
// array of Kubernetes runtime Object matched using the query from
// context. It utilizes a ResourceManager to execute the logic.
//
//	@Summary		SearchForResource returns an array of Kubernetes runtime Object matched using the query from context.
//	@Description	This endpoint returns an array of Kubernetes runtime Object matched using the query from context.
//	@Tags			resource
//	@Produce		json
//	@Success		200	{array}		runtime.Object	"Array of runtime.Object"
//	@Failure		400	{string}	string			"Bad Request"
//	@Failure		401	{string}	string			"Unauthorized"
//	@Failure		404	{string}	string			"Not Found"
//	@Failure		405	{string}	string			"Method Not Allowed"
//	@Failure		429	{string}	string			"Too Many Requests"
//	@Failure		500	{string}	string			"Internal Server Error"
//	@Router			/api/v1/resource/search [get]
func SearchForResource(resourceMgr *resource.ResourceManager, c *registry.ExtraConfig) http.HandlerFunc {
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
