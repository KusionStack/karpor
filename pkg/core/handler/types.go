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

// SuccessMessage is the default success message for successful responses.
const SuccessMessage = "OK"

// response defines the structure for API response payloads.
type response struct {
	Success   bool       `json:"success" yaml:"success"`                         // Indicates success status.
	Message   string     `json:"message" yaml:"message"`                         // Descriptive message.
	Data      any        `json:"data,omitempty" yaml:"data,omitempty"`           // Data payload.
	TraceID   string     `json:"traceID,omitempty" yaml:"traceID,omitempty"`     // Trace identifier.
	StartTime *time.Time `json:"startTime,omitempty" yaml:"startTime,omitempty"` // Request start time.
	EndTime   *time.Time `json:"endTime,omitempty" yaml:"endTime,omitempty"`     // Request end time.
	CostTime  Duration   `json:"costTime,omitempty" yaml:"costTime,omitempty"`   // Time taken for the request.
}

// Render is a no-op method that satisfies the render.Renderer interface.
func (rep *response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Response creates a standard API response renderer.
func Response(ctx context.Context, data any, err error) render.Renderer {
	resp := &response{}

	// Set the Success and Message fields based on the error parameter.
	if err == nil {
		resp.Success = true
		resp.Message = SuccessMessage
		resp.Data = data
	} else {
		resp.Success = false
		resp.Message = err.Error()
	}

	// Include the request trace ID if available.
	if requestID := middleware.GetReqID(ctx); len(requestID) > 0 {
		resp.TraceID = requestID
	}

	// Calculate and include timing details if a start time is set.
	if startTime := appmiddleware.GetStartTime(ctx); !startTime.IsZero() {
		endTime := time.Now()
		resp.StartTime = &startTime
		resp.EndTime = &endTime
		resp.CostTime = Duration(endTime.Sub(startTime))
	}

	return resp
}

// FailureResponse creates a response renderer for a failed request.
func FailureResponse(ctx context.Context, err error) render.Renderer {
	return Response(ctx, nil, err)
}

// SuccessResponse creates a response renderer for a successful request.
func SuccessResponse(ctx context.Context, data any) render.Renderer {
	return Response(ctx, data, nil)
}

// Duration is a custom type that represents a duration of time.
type Duration time.Duration

// MarshalJSON customizes JSON representation of the Duration type.
func (d Duration) MarshalJSON() (b []byte, err error) {
	// Format the duration as a string.
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(d).String())), nil
}
