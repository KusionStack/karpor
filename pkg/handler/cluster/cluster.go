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

package cluster

import (
	"encoding/json"
	"net/http"

	"github.com/KusionStack/karbour/pkg/controller/cluster"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/go-chi/chi/v5"
	"k8s.io/apiserver/pkg/server"
)

func Get(clusterCtrl *cluster.ClusterController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")
		clusterUnstructured, _ := clusterCtrl.GetCluster(r.Context(), client, cluster)
		result, _ := json.MarshalIndent(clusterUnstructured, "", "  ")
		w.Write(result)
	}
}

func GetYAML(clusterCtrl *cluster.ClusterController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")
		result, _ := clusterCtrl.GetYAMLForCluster(r.Context(), client, cluster)
		w.Write(result)
	}
}

func GetTopology(clusterCtrl *cluster.ClusterController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, cluster)
		topologyMap, _ := clusterCtrl.GetTopologyForCluster(r.Context(), client, cluster)
		result, _ := json.MarshalIndent(topologyMap, "", "  ")
		w.Write(result)
	}
}

func GetDetail(clusterCtrl *cluster.ClusterController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, cluster)
		clusterDetail, _ := clusterCtrl.GetDetailsForCluster(r.Context(), client, cluster)
		result, _ := json.MarshalIndent(clusterDetail, "", "  ")
		w.Write(result)
	}
}

func GetNamespace(clusterCtrl *cluster.ClusterController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		namespace := chi.URLParam(r, "namespaceName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, cluster)
		namespaceObj, _ := clusterCtrl.GetNamespaceForCluster(r.Context(), client, cluster, namespace)
		result, _ := json.MarshalIndent(namespaceObj, "", "  ")
		w.Write(result)
	}
}

func GetNamespaceTopology(clusterCtrl *cluster.ClusterController, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		namespace := chi.URLParam(r, "namespaceName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, cluster)
		topologyMap, _ := clusterCtrl.GetTopologyForClusterNamespace(r.Context(), client, cluster, namespace)
		result, _ := json.MarshalIndent(topologyMap, "", "  ")
		w.Write(result)
	}
}
