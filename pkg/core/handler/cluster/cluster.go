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
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/cluster"
	"github.com/KusionStack/karbour/pkg/multicluster"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"k8s.io/apiserver/pkg/server"
)

func Get(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")
		clusterUnstructured, _ := clusterMgr.GetCluster(r.Context(), client, cluster)
		result, _ := json.MarshalIndent(clusterUnstructured, "", "  ")
		w.Write(result)
	}
}

func GetYAML(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")
		result, _ := clusterMgr.GetYAMLForCluster(r.Context(), client, cluster)
		w.Write(result)
	}
}

func GetTopology(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, cluster)
		topologyMap, _ := clusterMgr.GetTopologyForCluster(r.Context(), client, cluster)
		result, _ := json.MarshalIndent(topologyMap, "", "  ")
		w.Write(result)
	}
}

func GetDetail(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, cluster)
		clusterDetail, _ := clusterMgr.GetDetailsForCluster(r.Context(), client, cluster)
		result, _ := json.MarshalIndent(clusterDetail, "", "  ")
		w.Write(result)
	}
}

func GetNamespace(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		namespace := chi.URLParam(r, "namespaceName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, cluster)
		namespaceObj, _ := clusterMgr.GetNamespaceForCluster(r.Context(), client, cluster, namespace)
		result, _ := json.MarshalIndent(namespaceObj, "", "  ")
		w.Write(result)
	}
}

func GetNamespaceTopology(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := chi.URLParam(r, "clusterName")
		namespace := chi.URLParam(r, "namespaceName")
		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, cluster)
		topologyMap, _ := clusterMgr.GetTopologyForClusterNamespace(r.Context(), client, cluster, namespace)
		result, _ := json.MarshalIndent(topologyMap, "", "  ")
		w.Write(result)
	}
}

// @Summary      Upload kubeConfig file for cluster
// @Description  Uploads a KubeConfig file for cluster, with a maximum size of 2MB, and valid file format is empty extension or JSON or YAML file.
// @Tags         cluster
// @Accept       multipart/form-data
// @Produce      plain
// @Param        file  formData  file    true  "Upload file with field name 'file'"
// @Success      200   {string}  string  "Returns the content of the uploaded KubeConfig file."
// @Failure      400   {string}  string  "The uploaded file is too large or the request is invalid."
// @Failure      500   {string}  string  "Internal server error."
// @Router       /api/v1/cluster/config/file [post]
func UpdateKubeConfig(w http.ResponseWriter, r *http.Request) {
	// Extract the context and logger from the request.
	ctx := r.Context()
	log := ctxutil.GetLogger(ctx)

	// Begin the auditing process, logging the start.
	log.Info("Starting get uploaded kubeconfig file in handler...")

	// Limit the size of the request body to prevent overflow.
	const maxUploadSize = 2 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "failed to parse multipart form")))
		return
	}

	// Retrieve the file from the parsed multipart form.
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "invalid request")))
		return
	}
	defer file.Close()

	// Check the file extension.
	log.Info("Uploaded file name:", "filename", fileHeader.Filename)
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != "" && ext != ".json" && ext != ".yaml" && ext != ".yml" {
		render.Render(w, r, handler.FailureResponse(ctx, errors.New("invalid file format, only empty extension and JSON and YAML files are allowed.")))
		return
	}

	// Read the contents of the file.
	buf := make([]byte, maxUploadSize)
	fileSize, err := file.Read(buf)
	if err != nil && err != io.EOF {
		render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "error reading file")))
		return
	}

	// Convert the bytes read to a string and return as response.
	fileContent := string(buf[:fileSize])
	render.JSON(w, r, handler.SuccessResponse(ctx, fileContent))
}
