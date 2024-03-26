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
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/cluster"
	"github.com/KusionStack/karbour/pkg/infra/multicluster"
	"github.com/KusionStack/karbour/pkg/util/clusterinstall"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pkg/errors"

	_ "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/tools/clientcmd"

	k8syaml "sigs.k8s.io/yaml"
)

// Get returns an HTTP handler function that reads a cluster
// detail. It utilizes a ClusterManager to execute the logic.
//
// @Summary      Get returns a cluster resource by name.
// @Description  This endpoint returns a cluster resource by name.
// @Tags         cluster
// @Produce      json
// @Param        format  query     string                     false  "The format of the response. Either in json or yaml"
// @Success      200     {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400     {string}  string                     "Bad Request"
// @Failure      401     {string}  string                     "Unauthorized"
// @Failure      404     {string}  string                     "Not Found"
// @Failure      405     {string}  string                     "Method Not Allowed"
// @Failure      429     {string}  string                     "Too Many Requests"
// @Failure      500     {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/cluster/{clusterName} [get]
func Get(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		cluster := chi.URLParam(r, "clusterName")
		logger.Info("Getting cluster...", "cluster", cluster)
		outputFormat := r.URL.Query().Get("format")

		client, err := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")
		if err != nil {
			render.Render(w, r, handler.FailureResponse(r.Context(), err))
			return
		}

		if strings.ToLower(outputFormat) == "yaml" {
			clusterYAML, err := clusterMgr.GetYAMLForCluster(r.Context(), client, cluster)
			handler.HandleResult(w, r, ctx, err, string(clusterYAML))
		} else {
			clusterUnstructured, err := clusterMgr.GetCluster(r.Context(), client, cluster)
			handler.HandleResult(w, r, ctx, err, clusterUnstructured)
		}
	}
}

// Create returns an HTTP handler function that creates a cluster
// resource. It utilizes a ClusterManager to execute the logic.
//
// @Summary      Create creates a cluster resource.
// @Description  This endpoint creates a new cluster resource using the payload.
// @Tags         cluster
// @Accept       plain
// @Accept       json
// @Produce      json
// @Param        request  body      ClusterPayload             true  "cluster to create (either plain text or JSON format)"
// @Success      200      {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400      {string}  string                     "Bad Request"
// @Failure      401      {string}  string                     "Unauthorized"
// @Failure      404      {string}  string                     "Not Found"
// @Failure      405      {string}  string                     "Method Not Allowed"
// @Failure      429      {string}  string                     "Too Many Requests"
// @Failure      500      {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/cluster/{clusterName} [post]
func Create(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		cluster := chi.URLParam(r, "clusterName")
		logger.Info("Creating cluster...", "cluster", cluster)

		// Decode the request body into the payload.
		payload := &ClusterPayload{}
		if err := payload.Decode(r); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")
		clusterCreated, err := clusterMgr.CreateCluster(r.Context(), client, cluster, payload.ClusterDisplayName, payload.ClusterDescription, payload.ClusterKubeConfig)
		handler.HandleResult(w, r, ctx, err, clusterCreated)
	}
}

// Update returns an HTTP handler function that updates a cluster
// resource. It utilizes a ClusterManager to execute the logic.
//
// @Summary      Update updates the cluster metadata by name.
// @Description  This endpoint updates the display name and description of an existing cluster resource.
// @Tags         cluster
// @Accept       plain
// @Accept       json
// @Produce      json
// @Param        request  body      ClusterPayload             true  "cluster to update (either plain text or JSON format)"
// @Success      200      {object}  unstructured.Unstructured  "Unstructured object"
// @Failure      400      {string}  string                     "Bad Request"
// @Failure      401      {string}  string                     "Unauthorized"
// @Failure      404      {string}  string                     "Not Found"
// @Failure      405      {string}  string                     "Method Not Allowed"
// @Failure      429      {string}  string                     "Too Many Requests"
// @Failure      500      {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/cluster/{clusterName}  [put]
func Update(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		cluster := chi.URLParam(r, "clusterName")
		logger.Info("Updating cluster metadata...", "cluster", cluster)

		// Decode the request body into the payload.
		payload := &ClusterPayload{}
		if err := payload.Decode(r); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")
		if payload.ClusterKubeConfig != "" {
			clusterUpdated, err := clusterMgr.UpdateCredential(r.Context(), client, cluster, payload.ClusterKubeConfig)
			handler.HandleResult(w, r, ctx, err, clusterUpdated)
		} else {
			clusterUpdated, err := clusterMgr.UpdateMetadata(r.Context(), client, cluster, payload.ClusterDisplayName, payload.ClusterDescription)
			handler.HandleResult(w, r, ctx, err, clusterUpdated)
		}
	}
}

// List returns an HTTP handler function that lists all cluster
// resources. It utilizes a ClusterManager to execute the logic.
//
// @Summary      List lists all cluster resources.
// @Description  This endpoint lists all cluster resources.
// @Tags         cluster
// @Produce      json
// @Param        summary     query     bool                       false  "Whether to display summary or not. Default to false"
// @Param        orderBy     query     string                     false  "The order to list the cluster. Default to order by name"
// @Param        descending  query     bool                       false  "Whether to sort the list in descending order. Default to false"
// @Success      200         {array}   unstructured.Unstructured  "List of cluster objects"
// @Failure      400         {string}  string                     "Bad Request"
// @Failure      401         {string}  string                     "Unauthorized"
// @Failure      404         {string}  string                     "Not Found"
// @Failure      405         {string}  string                     "Method Not Allowed"
// @Failure      429         {string}  string                     "Too Many Requests"
// @Failure      500         {string}  string                     "Internal Server Error"
// @Router       /rest-api/v1/clusters [get]
func List(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		logger.Info("Listing clusters...")

		orderBy := r.URL.Query().Get("orderBy")
		descending, _ := strconv.ParseBool(r.URL.Query().Get("descending"))
		summary, _ := strconv.ParseBool(r.URL.Query().Get("summary"))

		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")

		if summary {
			clusterSummary, err := clusterMgr.CountCluster(r.Context(), client, c.LoopbackClientConfig)
			handler.HandleResult(w, r, ctx, err, clusterSummary)
		} else {
			criteria, ok := sortCriteriaMap[orderBy]
			if !ok {
				criteria = cluster.ByName
			}
			clusterList, err := clusterMgr.ListCluster(r.Context(), client, criteria, descending)
			handler.HandleResult(w, r, ctx, err, clusterList)
		}
	}
}

// Delete returns an HTTP handler function that deletes a cluster
// resource. It utilizes a ClusterManager to execute the logic.
//
// @Summary      Delete removes a cluster resource by name.
// @Description  This endpoint deletes the cluster resource by name.
// @Tags         cluster
// @Produce      json
// @Success      200  {string}  string  "Operation status"
// @Failure      400  {string}  string  "Bad Request"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      404  {string}  string  "Not Found"
// @Failure      405  {string}  string  "Method Not Allowed"
// @Failure      429  {string}  string  "Too Many Requests"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /rest-api/v1/cluster/{clusterName} [delete]
func Delete(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)
		cluster := chi.URLParam(r, "clusterName")
		logger.Info("Deleting cluster...", "cluster", cluster)

		client, _ := multicluster.BuildMultiClusterClient(r.Context(), c.LoopbackClientConfig, "")
		err := clusterMgr.DeleteCluster(r.Context(), client, cluster)
		handler.HandleResult(w, r, ctx, err, "Cluster deleted")
	}
}

// @Summary      Upload kubeConfig file for cluster
// @Description  Uploads a KubeConfig file for cluster, with a maximum size of 2MB.
// @Tags         cluster
// @Accept       multipart/form-data
// @Produce      plain
// @Param        file         formData  file        true  "Upload file with field name 'file'"
// @Param        name         formData  string      true  "cluster name"
// @Param        displayName  formData  string      true  "cluster display name"
// @Param        description  formData  string      true  "cluster description"
// @Success      200          {object}  UploadData  "Returns the content of the uploaded KubeConfig file."
// @Failure      400          {string}  string      "The uploaded file is too large or the request is invalid."
// @Failure      500          {string}  string      "Internal server error."
// @Router       /rest-api/v1/cluster/config/file [post]
func UploadKubeConfig(clusterMgr *cluster.ClusterManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)

		// Begin the update process, logging the start.
		log.Info("Starting get uploaded kubeconfig file in handler...")

		// Limit the size of the request body to prevent overflow.
		const maxUploadSize = 2 << 20
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "failed to parse multipart form")))
			return
		}

		// Retrieve the file from the parsed multipart form.
		name := r.FormValue("name")
		displayName := r.FormValue("displayName")
		description := r.FormValue("description")
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "invalid request")))
			return
		}
		defer file.Close()

		log.Info("Uploaded filename", "filename", fileHeader.Filename)

		// Read the contents of the file.
		buf := make([]byte, maxUploadSize)
		fileSize, err := file.Read(buf)
		if err != nil && errors.Is(err, io.EOF) {
			render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "error reading file")))
			return
		}
		plainContent := string(buf[:fileSize])

		// Create new restConfig from uploaded kubeconfig.
		restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(plainContent))
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "error create new rest config from uploaded kubeconfig")))
			return
		}

		// Convert the rest.Config to Cluster object.
		clusterObj, err := clusterinstall.ConvertKubeconfigToCluster(name, displayName, description, restConfig)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "error convert kubeconfig to cluster")))
			return
		}
		unstructuredClusterMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(clusterObj)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "error convert cluster to unstructured obj")))
			return
		}
		clusterUnstructured := &unstructured.Unstructured{}
		clusterUnstructured.SetUnstructuredContent(unstructuredClusterMap)

		// Sanitize the cluster object.
		sanitizedUnstructuredClusterObj, err := cluster.SanitizeUnstructuredCluster(ctx, clusterUnstructured)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "error sanitize unstructured obj")))
			return
		}

		clusterYAML, err := k8syaml.Marshal(sanitizedUnstructuredClusterObj)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, errors.Wrapf(err, "error marshal unstructured obj")))
			return
		}

		// Convert the bytes read to a string and return as response.
		data := &UploadData{
			FileName:                fileHeader.Filename,
			FileSize:                fileSize,
			Content:                 plainContent,
			SanitizedClusterContent: string(clusterYAML),
		}
		render.JSON(w, r, handler.SuccessResponse(ctx, data))
	}
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
// @Router       /rest-api/v1/cluster/config/validate [post]
func ValidateKubeConfig(clusterMgr *cluster.ClusterManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)

		// Begin the auditing process, logging the start.
		log.Info("Starting validate kubeconfig in handler...")

		// Decode the request body into the payload.
		payload := &ValidatePayload{}
		if err := payload.Decode(r); err != nil {
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
