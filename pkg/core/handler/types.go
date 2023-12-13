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

package handler

import (
	"fmt"
	"net/http"
	"time"
)

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

// Payload is an interface for incoming requests payloads
// Each handler should implement this interface to parse payloads
type Payload interface {
	Decode(*http.Request) error // Decode returns the payload object with the decoded
}

// Render is a no-op method that satisfies the render.Renderer interface.
func (rep *response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Duration is a custom type that represents a duration of time.
type Duration time.Duration

// MarshalJSON customizes JSON representation of the Duration type.
func (d Duration) MarshalJSON() (b []byte, err error) {
	// Format the duration as a string.
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(d).String())), nil
}
