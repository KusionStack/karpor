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
// @Description  Uploads a KubeConfig file for cluster, with a maximum size of 2MB, and the valid file extension is "", ".yaml", ".yml", ".json", ".kubeconfig", ".kubeconf".
// @Tags         cluster
// @Accept       multipart/form-data
// @Produce      plain
// @Param        file  formData  file        true  "Upload file with field name 'file'"
// @Success      200   {object}  UploadData  "Returns the content of the uploaded KubeConfig file."
// @Failure      400   {string}  string      "The uploaded file is too large or the request is invalid."
// @Failure      500   {string}  string      "Internal server error."
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
	log.Info("Uploaded filename", "filename", fileHeader.Filename)
	if !isAllowedExtension(fileHeader.Filename) {
		render.Render(w, r, handler.FailureResponse(
			ctx, errors.New("invalid file format, only '', .yaml, .yml, .json, .kubeconfig, .kubeconf are allowed.")))
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
	data := &UploadData{
		FileName: fileHeader.Filename,
		Content:  string(buf[:fileSize]),
		FileSize: fileSize,
	}
	render.JSON(w, r, handler.SuccessResponse(ctx, data))
}

// ValidateKubeConfig returns an HTTP handler function to validate a KubeConfig.
//
// @Summary      Validate KubeConfig
// @Description  Validates the provided KubeConfig using cluster manager methods.
// @Tags         cluster
// @Accept       plain
// @Accept       json
// @Produce      json
// @Param        request  body      ValidatePayload  true  "KubeConfig payload to validate"
// @Success      200      {string}  string           "Verification passed server version"
// @Failure      400      {object}  string           "Bad Request"
// @Failure      401      {object}  string           "Unauthorized"
// @Failure      429      {object}  string           "Too Many Requests"
// @Failure      404      {object}  string           "Not Found"
// @Failure      500      {object}  string           "Internal Server Error"
// @Router       /api/v1/cluster/config/validate [post]
func ValidateKubeConfig(clusterMgr *cluster.ClusterManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)

		// Begin the auditing process, logging the start.
		log.Info("Starting validate kubeconfig in handler...")

		// Decode the request body into the payload.
		payload := &ValidatePayload{}
		if err := decode(r, payload); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Log successful decoding of the request body.
		sanitizeConfig, _ := clusterMgr.SanitizeKubeConfigWithYAML(ctx, payload.KubeConfig)
		log.Info("Successfully decoded the request body to payload, and sanitize the kubeconfig in payload",
			"sanitizeKubeConfig", sanitizeConfig)

		// Validate the specified kube config.
		if info, err := clusterMgr.ValidateKubeConfigWithYAML(ctx, payload.KubeConfig); err == nil {
			render.JSON(w, r, handler.SuccessResponse(ctx, info))
		} else {
			render.Render(w, r, handler.FailureResponse(ctx, err))
		}
	}
}

// isAllowedExtension checks if the provided file name has a permitted extension.
func isAllowedExtension(filename string) bool {
	allowedExtensions := []string{"", ".yaml", ".yml", ".json", ".kubeconfig", ".kubeconf"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
