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

package audit

import (
	"io"
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/audit"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
)

// Audit returns an HTTP handler function that performs auditing of manifest
// data. It utilizes an AuditManager to execute the audit logic.
// @Summary AuditHandler audits the provided manifest.
// @Description This endpoint audits the provided manifest for issues.
// @Tags audit
// @Accept plain, json
// @Produce json
// @Param manifest body string true "Manifest data to audit (either plain text or JSON format)"
// @Success 200 {object} AuditResponse "Audit results"
// @Failure 400 {object} FailureResponse "Error details if manifest cannot be processed"
// @Router /audit [post]
func Audit(auditMgr *audit.AuditManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)

		// Begin the auditing process, logging the start.
		log.Info("Starting audit of the specified manifest in handler ...")

		// Initialize an empty payload to hold the manifest data.
		payload := &Payload{}

		// Check if the content type is plain text, read it as such.
		if render.GetRequestContentType(r) == render.ContentTypePlainText {
			// Read the request body.
			body, err := io.ReadAll(r.Body)
			defer r.Body.Close() // Ensure the body is closed after reading.
			if err != nil {
				// Handle any reading errors by sending a failure response.
				render.Render(w, r, handler.FailureResponse(ctx, err))
				return
			}
			// Set the read content as the manifest payload.
			payload.Manifest = string(body)
		} else {
			// For non-plain text, decode the JSON body into the payload.
			if err := render.DecodeJSON(r.Body, payload); err != nil {
				// Handle JSON decoding errors.
				render.Render(w, r, handler.FailureResponse(ctx, err))
				return
			}
		}

		// Log successful decoding of the request body.
		log.Info("Successfully decoded the request body to payload",
			"payload", payload)

		// Perform the audit using the manager and the provided manifest.
		issues, err := auditMgr.Audit(ctx, payload.Manifest)
		if err != nil {
			// Handle audit errors by sending a failure response.
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Send a success response with the audit issues.
		render.JSON(w, r, handler.SuccessResponse(ctx, issues))
	}
}
