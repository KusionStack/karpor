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
	"fmt"
	"net/http"
	"net/url"

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/go-chi/chi/v5"
)

func BuildResourceFromParam(r *http.Request) *Resource {
	apiVersion := chi.URLParam(r, "apiVersion")
	if r.URL.RawPath != "" {
		apiVersion, _ = url.PathUnescape(apiVersion)
	}
	res := Resource{
		Cluster:    chi.URLParam(r, "clusterName"),
		APIVersion: apiVersion,
		Kind:       chi.URLParam(r, "kind"),
		Namespace:  chi.URLParam(r, "namespaceName"),
		Name:       chi.URLParam(r, "resourceName"),
	}
	return &res
}

func BuildLocatorFromQuery(r *http.Request) (*core.Locator, error) {
	apiVersion := r.URL.Query().Get("apiVersion")
	if r.URL.RawPath != "" {
		apiVersion, _ = url.PathUnescape(apiVersion)
	}
	cluster := r.URL.Query().Get("cluster")
	if cluster == "" {
		return &core.Locator{}, fmt.Errorf("cluster cannot be empty")
	}
	res := core.Locator{
		Cluster:    cluster,
		APIVersion: apiVersion,
		Kind:       r.URL.Query().Get("kind"),
		Namespace:  r.URL.Query().Get("namespace"),
		Name:       r.URL.Query().Get("name"),
	}
	return &res, nil
}
