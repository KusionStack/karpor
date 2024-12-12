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
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/kubernetes"

	"github.com/KusionStack/karpor/pkg/core/manager/cluster"
	"github.com/KusionStack/karpor/pkg/infra/multicluster"
	"github.com/KusionStack/karpor/pkg/util/ctxutil"
)

const TimeFormat = "2006-01-02T15:04:05Z"

// Event represents a Kubernetes event with additional fields
type Event struct {
	Type           string `json:"type"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	Count          int32  `json:"count"`
	LastTimestamp  string `json:"lastTimestamp"`
	FirstTimestamp string `json:"firstTimestamp"`
}

// GetEvents returns an HTTP handler function that streams events for a resource using SSE
//
// @Summary      Stream resource events using Server-Sent Events
// @Description  This endpoint streams resource events in real-time using SSE. It supports event type filtering and automatic updates.
// @Tags         insight
// @Produce      text/event-stream
// @Param        cluster     path      string  true   "The cluster name"
// @Param        namespace   path      string  true   "The namespace name"
// @Param        name        path      string  true   "The resource name"
// @Param        kind        query     string  true   "The resource kind"
// @Param        apiVersion  query     string  true   "The resource API version"
// @Param        type        query     string  false  "Event type filter (Normal or Warning)"
// @Success      200         {array}   Event
// @Failure      400         {string}  string  "Bad Request"
// @Failure      401         {string}  string  "Unauthorized"
// @Failure      404         {string}  string  "Not Found"
// @Router       /insight/aggregator/event/{cluster}/{namespace}/{name} [get]
func GetEvents(clusterMgr *cluster.ClusterManager, c *server.CompletedConfig) http.HandlerFunc {
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
		kind := r.URL.Query().Get("kind")
		apiVersion := r.URL.Query().Get("apiVersion")
		eventType := r.URL.Query().Get("type")

		if cluster == "" || namespace == "" || name == "" || kind == "" || apiVersion == "" {
			writeEventSSEError(w, "missing required parameters")
			return
		}

		// Build multi-cluster client
		client, err := multicluster.BuildMultiClusterClient(ctx, c.LoopbackClientConfig, cluster)
		if err != nil {
			writeEventSSEError(w, "failed to build multi-cluster client: "+err.Error())
			return
		}
		// Get single cluster clientset
		clusterClient := client.ClientSet

		logger.Info("Streaming resource events...", "cluster", cluster, "namespace", namespace, "name", name, "kind", kind, "apiVersion", apiVersion, "type", eventType)

		// Create a ticker to periodically fetch events
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// Create a done channel to handle client disconnection
		done := make(chan bool)
		go func() {
			<-ctx.Done()
			done <- true
		}()

		// Send initial events
		if err := streamEvents(ctx, w, clusterClient, namespace, name, kind, apiVersion, eventType); err != nil {
			writeEventSSEError(w, "failed to get events: "+err.Error())
			return
		}

		// Keep streaming events until client disconnects
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := streamEvents(ctx, w, clusterClient, namespace, name, kind, apiVersion, eventType); err != nil {
					writeEventSSEError(w, "failed to get events: "+err.Error())
					return
				}
			}
		}
	}
}

func streamEvents(ctx context.Context, w http.ResponseWriter, client kubernetes.Interface, namespace, name, kind, apiVersion, eventType string) error {
	events, err := getResourceEvents(ctx, client, namespace, name, kind, apiVersion, eventType)
	if err != nil {
		return err
	}

	data, err := json.Marshal(events)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("data: " + string(data) + "\n\n"))
	if err != nil {
		return err
	}

	w.(http.Flusher).Flush()
	return nil
}

func getResourceEvents(ctx context.Context, client kubernetes.Interface, namespace, name, kind, apiVersion, eventType string) ([]Event, error) {
	fieldSelector := "involvedObject.name=" + name +
		",involvedObject.namespace=" + namespace +
		",involvedObject.kind=" + kind +
		",involvedObject.apiVersion=" + apiVersion

	if eventType != "" {
		fieldSelector += ",type=" + eventType
	}

	k8sEvents, err := client.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0, len(k8sEvents.Items))
	for _, e := range k8sEvents.Items {
		events = append(events, Event{
			Type:           e.Type,
			Reason:         e.Reason,
			Message:        e.Message,
			Count:          e.Count,
			LastTimestamp:  e.LastTimestamp.UTC().Format(TimeFormat),
			FirstTimestamp: e.FirstTimestamp.UTC().Format(TimeFormat),
		})
	}

	// Sort events by last timestamp in descending order (same as kubectl describe)
	sort.Slice(events, func(i, j int) bool {
		iTime, _ := time.Parse(time.RFC3339, events[i].LastTimestamp)
		jTime, _ := time.Parse(time.RFC3339, events[j].LastTimestamp)
		return iTime.After(jTime)
	})

	return events, nil
}

// writeEventSSEError writes an error message to the SSE stream using Event format
func writeEventSSEError(w http.ResponseWriter, errMsg string) {
	event := Event{
		Type:          "Warning",
		Reason:        "Error",
		Message:       errMsg,
		Count:         1,
		LastTimestamp: time.Now().Format(TimeFormat),
	}
	data, _ := json.Marshal([]Event{event})
	w.Write([]byte("data: " + string(data) + "\n\n"))
	w.(http.Flusher).Flush()
}
