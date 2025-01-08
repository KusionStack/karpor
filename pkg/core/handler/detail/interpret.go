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

package detail

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KusionStack/karpor/pkg/core/handler"
	"github.com/KusionStack/karpor/pkg/core/manager/ai"
	"github.com/KusionStack/karpor/pkg/util/ctxutil"
	"k8s.io/apiserver/pkg/server"
)

// InterpretRequest represents the request body for YAML interpretation
type InterpretRequest struct {
	YAML     string `json:"yaml"`
	Language string `json:"language"`
}

// InterpretYAML returns an HTTP handler function that performs AI interpretation on YAML content
//
// @Summary      Interpret YAML using AI
// @Description  This endpoint analyzes YAML content using AI to provide detailed interpretation and insights
// @Tags         insight
// @Accept       json
// @Produce      text/event-stream
// @Param        request  body      InterpretRequest  true  "The YAML content to interpret"
// @Success      200      {object}  ai.InterpretEvent
// @Failure      400      {string}  string  "Bad Request"
// @Failure      401      {string}  string  "Unauthorized"
// @Failure      429      {string}  string  "Too Many Requests"
// @Failure      404      {string}  string  "Not Found"
// @Failure      500      {string}  string  "Internal Server Error"
// @Router       /insight/yaml/interpret/stream [post]
func InterpretYAML(aiMgr *ai.AIManager, c *server.CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request
		ctx := r.Context()
		logger := ctxutil.GetLogger(ctx)

		// Begin the interpretation process, logging the start
		logger.Info("Starting YAML interpretation in handler ...")

		// Parse request body
		var req InterpretRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handler.FailureRender(ctx, w, r, fmt.Errorf("invalid request format: %v", err))
			return
		}

		// Validate request
		if req.YAML == "" {
			handler.FailureRender(ctx, w, r, fmt.Errorf("YAML content is required"))
			return
		}
		if req.Language == "" {
			req.Language = "English" // Default to English if language not specified
		}

		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Accel-Buffering", "no")

		flusher, ok := w.(http.Flusher)
		if !ok {
			handler.FailureRender(ctx, w, r, fmt.Errorf("streaming unsupported"))
			return
		}

		// Create channel for interpretation events
		eventChan := make(chan *ai.InterpretEvent, 10)
		go func() {
			if err := aiMgr.InterpretYAML(ctx, req.YAML, req.Language, eventChan); err != nil {
				logger.Error(err, "Failed to interpret YAML")
				// Error will be sent through eventChan
			}
		}()

		// Stream events to client
		for event := range eventChan {
			data, err := json.Marshal(event)
			if err != nil {
				logger.Error(err, "Failed to marshal event")
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
