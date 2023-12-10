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
