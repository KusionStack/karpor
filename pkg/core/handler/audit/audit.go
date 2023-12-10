package audit

import (
	"io"
	"net/http"

	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/core/manager/audit"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
	"github.com/go-chi/render"
)

func Audit(auditMgr *audit.AuditManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := ctxutil.GetLogger(ctx)

		log.Info("Starting audit the specified manifest in handler ...")

		payload := &Payload{}
		if render.GetRequestContentType(r) == render.ContentTypePlainText {
			body, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				render.Render(w, r, handler.FailureResponse(ctx, err))
				return
			}
			payload.Manifest = string(body)
		} else {
			if err := render.DecodeJSON(r.Body, payload); err != nil {
				render.Render(w, r, handler.FailureResponse(ctx, err))
				return
			}
		}

		log.Info("Successfully decoded the request body to payload", "payload", payload)

		issues, err := auditMgr.Audit(ctx, payload.Manifest)
		if err != nil {
			render.Render(w, r, handler.FailureResponse(ctx, err))
			return
		}

		render.JSON(w, r, handler.SuccessResponse(ctx, issues))
	}
}
