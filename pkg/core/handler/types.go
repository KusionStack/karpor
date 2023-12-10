package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	appmiddleware "github.com/KusionStack/karbour/pkg/middleware"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const SuccessMessage = "OK"

type response struct {
	Success   bool       `json:"success" yaml:"success"`
	Message   string     `json:"message" yaml:"message"`
	Data      any        `json:"data,omitempty" yaml:"data,omitempty"`
	TraceID   string     `json:"traceID,omitempty" yaml:"traceID,omitempty"`
	StartTime *time.Time `json:"startTime,omitempty" yaml:"startTime,omitempty"`
	EndTime   *time.Time `json:"endTime,omitempty" yaml:"endTime,omitempty"`
	CostTime  Duration   `json:"costTime,omitempty" yaml:"costTime,omitempty"`
}

func (rep *response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Response(ctx context.Context, data any, err error) render.Renderer {
	resp := &response{}

	if err == nil {
		resp.Success = true
		resp.Message = SuccessMessage
		resp.Data = data
	} else {
		resp.Success = false
		resp.Message = err.Error()
	}

	if requestID := middleware.GetReqID(ctx); len(requestID) > 0 {
		resp.TraceID = requestID
	}

	if startTime := appmiddleware.GetStartTime(ctx); !startTime.IsZero() {
		endTime := time.Now()
		resp.StartTime = &startTime
		resp.EndTime = &endTime
		resp.CostTime = Duration(endTime.Sub(startTime))
	}

	return resp
}

func FailureResponse(ctx context.Context, err error) render.Renderer {
	return Response(ctx, nil, err)
}

func SuccessResponse(ctx context.Context, data any) render.Renderer {
	return Response(ctx, data, nil)
}

type Duration time.Duration

func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, (time.Duration(d)).String())), nil
}
