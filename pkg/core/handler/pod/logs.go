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

package pod

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
			writeSSEError(w, "cluster, namespace and name are required")
			return
		}

		// Build multi-cluster client
		client, err := multicluster.BuildMultiClusterClient(ctx, c.LoopbackClientConfig, cluster)
		if err != nil {
			writeSSEError(w, fmt.Sprintf("failed to build multi-cluster client: %v", err))
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
			writeSSEError(w, fmt.Sprintf("failed to get pod logs: %v", err))
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
					writeSSEError(w, fmt.Sprintf("failed to marshal log entry: %v", err))
					return
				}

				fmt.Fprintf(w, "data: %s\n\n", data)
				w.(http.Flusher).Flush()
			}
		}

		if err := scanner.Err(); err != nil {
			writeSSEError(w, fmt.Sprintf("error reading log stream: %v", err))
		}
	}
}

// writeSSEError writes an error message to the SSE stream
func writeSSEError(w http.ResponseWriter, errMsg string) {
	logEntry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Error:     errMsg,
	}
	data, _ := json.Marshal(logEntry)
	fmt.Fprintf(w, "data: %s\n\n", data)
	w.(http.Flusher).Flush()
}
