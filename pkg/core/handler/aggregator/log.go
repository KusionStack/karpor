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
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/KusionStack/karpor/pkg/core/manager/ai"
	"github.com/KusionStack/karpor/pkg/core/manager/cluster"
	"github.com/KusionStack/karpor/pkg/infra/multicluster"
	"github.com/KusionStack/karpor/pkg/util/ctxutil"
	"github.com/go-chi/chi/v5"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/utils/pointer"
)

// LogEntry represents a single log entry with timestamp and content
type LogEntry struct {
	Timestamp string `json:"timestamp,omitempty"`
	Content   string `json:"content"`
	Error     string `json:"error,omitempty"`
}

// GetPodLogs returns an HTTP handler function that streams Pod logs using Server-Sent Events
//
// @Summary      Stream pod logs using Server-Sent Events
// @Description  This endpoint streams pod logs in real-time using SSE. It supports container selection and automatic reconnection.
// @Tags         insight
// @Produce      text/event-stream,application/json
// @Param        cluster     path      string  true   "The cluster name"
// @Param        namespace   path      string  true   "The namespace name"
// @Param        name        path      string  true   "The pod name"
// @Param        container   query     string  false  "The container name (optional if pod has only one container)"
// @Param        since       query     string  false  "Only return logs newer than a relative duration like 5s, 2m, or 3h"
// @Param        sinceTime   query     string  false  "Only return logs after a specific date (RFC3339)"
// @Param        timestamps  query     bool    false  "Include timestamps in log output"
// @Param        tailLines   query     int     false  "Number of lines from the end of the logs to show"
// @Param        download    query     bool    false  "Download logs as file instead of streaming"
// @Success      200         {object}  LogEntry
// @Failure      400         {string}  string  "Bad Request"
// @Failure      401         {string}  string  "Unauthorized"
// @Failure      404         {string}  string  "Not Found"
// @Router       /insight/aggregator/log/pod/{cluster}/{namespace}/{name} [get]
func GetPodLogs(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		// Get parameters from URL path and query
		cluster := chi.URLParam(r, "cluster")
		namespace := chi.URLParam(r, "namespace")
		name := chi.URLParam(r, "name")
		container := r.URL.Query().Get("container")
		download := r.URL.Query().Get("download") == "true"

		if cluster == "" || namespace == "" || name == "" {
			writeLogSSEError(w, "cluster, namespace and name are required")
			return
		}

		// Configure log streaming options
		opts, err := buildPodLogOptions(r.URL.Query(), download)
		if err != nil {
			writeLogSSEError(w, err.Error())
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

		// Get log stream from the pod
		req := clusterClient.CoreV1().Pods(namespace).GetLogs(name, opts)
		stream, err := req.Stream(ctx)
		if err != nil {
			writeLogSSEError(w, fmt.Sprintf("failed to get pod logs: %v", err))
			return
		}
		defer stream.Close()

		// Handle download request
		if download {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s.log", name, container))
			io.Copy(w, stream)
			return
		}

		// Set SSE headers for streaming
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Create a scanner to read the log stream
		scanner := bufio.NewScanner(stream)
		scanner.Split(bufio.ScanLines)

		// Send logs as SSE events
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 && opts.Timestamps {
				logEntry := LogEntry{
					Timestamp: parts[0],
					Content:   parts[1],
				}
				data, err := json.Marshal(logEntry)
				if err != nil {
					logger.Error(err, "Failed to marshal log entry")
					continue
				}
				fmt.Fprintf(w, "data: %s\n\n", data)
			} else {
				logEntry := LogEntry{
					Content: line,
				}
				data, err := json.Marshal(logEntry)
				if err != nil {
					logger.Error(err, "Failed to marshal log entry")
					continue
				}
				fmt.Fprintf(w, "data: %s\n\n", data)
			}
			w.(http.Flusher).Flush()
		}

		if err := scanner.Err(); err != nil {
			logger.Error(err, "Error reading log stream")
		}
	}
}

// buildPodLogOptions builds the pod log options from the query parameters of
// a log request. It returns an error if since, sinceTime or tailLines is set
// to a value that cannot be parsed, so that the filter is not silently
// ignored.
func buildPodLogOptions(query url.Values, download bool) (*corev1.PodLogOptions, error) {
	opts := &corev1.PodLogOptions{
		Container:  query.Get("container"),
		Follow:     !download, // Don't follow if downloading
		Timestamps: query.Get("timestamps") == "true",
	}

	// Parse and set since time
	if since := query.Get("since"); since != "" {
		duration, err := time.ParseDuration(since)
		if err != nil || duration <= 0 {
			return nil, fmt.Errorf("invalid since %q: must be a positive duration such as 30s, 5m or 2h", since)
		}
		// Round up so that sub-second durations are not sent as 0, which the
		// Kubernetes API rejects.
		opts.SinceSeconds = pointer.Int64(int64(math.Ceil(duration.Seconds())))
	} else if sinceTime := query.Get("sinceTime"); sinceTime != "" {
		t, err := time.Parse(time.RFC3339, sinceTime)
		if err != nil {
			return nil, fmt.Errorf("invalid sinceTime %q: must be an RFC3339 timestamp such as 2006-01-02T15:04:05Z", sinceTime)
		}
		opts.SinceTime = &metav1.Time{Time: t}
	}

	// Parse and set tail lines
	if tailLines := query.Get("tailLines"); tailLines != "" {
		lines, err := strconv.ParseInt(tailLines, 10, 64)
		if err != nil || lines < 0 {
			return nil, fmt.Errorf("invalid tailLines %q: must be a non-negative integer", tailLines)
		}
		opts.TailLines = pointer.Int64(lines)
	} else if !download {
		opts.TailLines = pointer.Int64(1000) // Default to last 1000 lines for streaming
	}

	return opts, nil
}

// writeLogSSEError writes an error message to the SSE stream
func writeLogSSEError(w http.ResponseWriter, errMsg string) {
	logEntry := LogEntry{
		Error: errMsg,
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
// @Param        request  body      DiagnoseRequest  true  "The logs to analyze"
// @Success      200      {object}  ai.DiagnosisEvent
// @Failure      400      {string}  string  "Bad Request"
// @Failure      500      {string}  string  "Internal Server Error"
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
