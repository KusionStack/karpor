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
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/audit"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
)

// Audit returns an HTTP handler function that performs auditing of manifest
// data. It utilizes an AuditManager to execute the audit logic.
//
//	@Summary		AuditHandler audits the provided manifest.
//	@Description	This endpoint audits the provided manifest for issues.
//	@Tags			audit
//	@Accept			plain
//	@Accept			json
//	@Produce		json
//	@Param			request	body		AuditPayload	true	"Manifest data to audit (either plain text or JSON format)"
//	@Success		200		{array}		scanner.Issue	"Audit results"
//	@Failure		400		{string}	string			"Bad Request"
//	@Failure		401		{string}	string			"Unauthorized"
//	@Failure		429		{string}	string			"Too Many Requests"
//	@Failure		404		{string}	string			"Not Found"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Router			/api/v1/audit [post]
func Audit(auditMgr *audit.AuditManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)
		// Begin the auditing process, logging the start.
		log.Info("Starting audit of the specified manifest in handler ...")

		// Decode the request body into the payload.
		payload := &AuditPayload{}
		if err := payload.Decode(r); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Log successful decoding of the request body.
		log.Info("Successfully decoded the request body to payload",
			"payload", payload)

		// Perform the audit using the manager and the provided manifest.
		issues, err := auditMgr.Audit(ctx, payload.Manifest)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		render.JSON(w, r, handler.SuccessResponse(ctx, issues))
	}
}

// Score returns an HTTP handler function that calculates a score for the
// audited manifest. It utilizes an AuditManager to compute the score based
// on detected issues.
//
//	@Summary		ScoreHandler calculates a score for the audited manifest.
//	@Description	This endpoint calculates a score for the provided manifest based on the number and severity of issues detected during the audit.
//	@Tags			audit
//	@Accept			text/plain
//	@Accept			application/json
//	@Produce		json
//	@Param			request	body		AuditPayload	true	"Manifest data to calculate score for (either plain text or JSON format)"
//	@Success		200		{object}	audit.ScoreData	"Score calculation result"
//	@Failure		400		{string}	string			"Bad Request"
//	@Failure		401		{string}	string			"Unauthorized"
//	@Failure		429		{string}	string			"Too Many Requests"
//	@Failure		404		{string}	string			"Not Found"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Router			/api/v1/audit/score [post]
func Score(auditMgr *audit.AuditManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the context and logger from the request.
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)

		// Begin the auditing process, logging the start.
		log.Info("Starting calculate score with specified manifest in handler...")

		// Decode the request body into the payload.
		payload := &AuditPayload{}
		if err := payload.Decode(r); err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Log successful decoding of the request body.
		log.Info("Successfully decoded the request body to payload",
			"payload", payload)

		// Perform the audit to gather issues for score calculation.
		issues, err := auditMgr.Audit(ctx, payload.Manifest)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		// Calculate score using the audit issues.
		data, err := auditMgr.Score(ctx, issues)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		render.JSON(w, r, handler.SuccessResponse(ctx, data))
	}
}
