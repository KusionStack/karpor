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

package aggregator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KusionStack/karpor/pkg/core/manager/ai"
	"github.com/KusionStack/karpor/pkg/core/manager/cluster"
	"github.com/KusionStack/karpor/pkg/infra/multicluster"
	"github.com/KusionStack/karpor/pkg/util/ctxutil"
	"github.com/go-chi/chi/v5"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/utils/pointer"
)

// LogEntry represents a single log entry with timestamp and content
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
	Error     string `json:"error,omitempty"`
}

// GetPodLogs returns an HTTP handler function that streams Pod logs using Server-Sent Events
//
// @Summary      Stream pod logs using Server-Sent Events
// @Description  This endpoint streams pod logs in real-time using SSE. It supports container selection and automatic reconnection.
// @Tags         insight
// @Produce      text/event-stream
// @Param        cluster    path      string  true   "The cluster name"
// @Param        namespace  path      string  true   "The namespace name"
// @Param        name       path      string  true   "The pod name"
// @Param        container  query     string  false  "The container name (optional if pod has only one container)"
// @Success      200        {object}  LogEntry
// @Failure      400        {string}  string  "Bad Request"
// @Failure      401        {string}  string  "Unauthorized"
// @Failure      404        {string}  string  "Not Found"
// @Router       /insight/aggregator/log/pod/{cluster}/{namespace}/{name} [get]
func GetPodLogs(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Extract the context and logger from the request
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		// Get parameters from URL path and query
		cluster := chi.URLParam(r, "cluster")
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")
		container := r.URL.Query().Get("container")

		if cluster == "" || namespace == "" || name == "" {
			writeLogSSEError(w, "cluster, namespace and name are required")
			return
		}

		// Build multi-cluster client
		client, err := multicluster.BuildMultiClusterClient(ctx, c.LoopbackClientConfig, cluster)
		if err != nil {
			writeLogSSEError(w, fmt.Sprintf("failed to build multi-cluster client: %v", err))
			return
		}
		// Get single cluster clientset
		clusterClient := client.ClientSet

		logger.Info("Getting pod logs...", "cluster", cluster, "namespace", namespace, "pod", name, "container", container)

		// Configure log streaming options
		opts := &corev1.PodLogOptions{
			Container: container,
			Follow:    true,
			TailLines: pointer.Int64(1000),
		}

		// Get log stream from the pod
		req := clusterClient.CoreV1().Pods(namespace).GetLogs(name, opts)
		stream, err := req.Stream(ctx)
		if err != nil {
			writeLogSSEError(w, fmt.Sprintf("failed to get pod logs: %v", err))
			return
		}
		defer stream.Close()

		// Create a done channel to handle client disconnection
		done := r.Context().Done()
		go func() {
			<-done
			stream.Close()
		}()

		// Read and send logs
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			select {
			case <-done:
				return
			default:
				logEntry := LogEntry{
					Timestamp: time.Now().Format(time.RFC3339Nano),
					Content:   scanner.Text(),
				}

				data, err := json.Marshal(logEntry)
				if err != nil {
					writeLogSSEError(w, fmt.Sprintf("failed to marshal log entry: %v", err))
					return
				}

				fmt.Fprintf(w, "data: %s\n\n", data)
				w.(http.Flusher).Flush()
			}
		}

		if err := scanner.Err(); err != nil {
			writeLogSSEError(w, fmt.Sprintf("error reading log stream: %v", err))
		}
	}
}

// writeLogSSEError writes an error message to the SSE stream
func writeLogSSEError(w http.ResponseWriter, errMsg string) {
	logEntry := LogEntry{
		Timestamp: time.Now().Format(TimeFormat),
		Error:     errMsg,
	}
	data, _ := json.Marshal(logEntry)
	fmt.Fprintf(w, "data: %s\n\n", data)
	w.(http.Flusher).Flush()
}

// DiagnoseRequest represents the request body for log diagnosis
type DiagnoseRequest struct {
	Logs     []string `json:"logs"`
	Language string   `json:"language"` // Language code for AI response
}

// DiagnoseResponse represents the response for log diagnosis
type DiagnoseResponse struct {
	Diagnosis string `json:"diagnosis"`
}

// DiagnosePodLogs returns an HTTP handler function that performs AI diagnosis on pod logs
//
// @Summary      Diagnose pod logs using AI
// @Description  This endpoint analyzes pod logs using AI to identify issues and provide solutions
// @Tags         insight
// @Accept       json
// @Produce      text/event-stream
// @Param        request body DiagnoseRequest true "The logs to analyze"
// @Success      200  {object}  DiagnosisEvent
// @Failure      400  {string}  string "Bad Request"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /insight/aggregator/log/diagnosis/stream [post]
func DiagnosePodLogs(aiMgr *ai.AIManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		if err := ai.CheckAIManager(aiMgr); err != nil {
			logger.Error(err, "AI manager is not available")
			http.Error(w, "AI service is not available", http.StatusServiceUnavailable)
			return
		}

		// Parse request body
		var req DiagnoseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Accel-Buffering", "no")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Create channel for diagnosis events
		eventChan := make(chan *ai.DiagnosisEvent, 10)
		go func() {
			if err := aiMgr.DiagnoseLogs(ctx, req.Logs, req.Language, eventChan); err != nil {
				// Error already sent through eventChan
				return
			}
		}()

		// Stream events to client
		for event := range eventChan {
			data, err := json.Marshal(event)
			if err != nil {
				logger.Error(err, "Failed to marshal diagnosis event")
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
