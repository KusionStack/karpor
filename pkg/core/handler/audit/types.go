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

	"github.com/KusionStack/karbour/pkg/core"
	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

// Ensure that ClusterPayload implements the handler.Payload interface.
var _ handler.Payload = &AuditManifestPayload{}

// Payload represents the structure for audit request data. It includes the
// manifest which is typically a string containing declarative configuration
// data that needs to be audited.
type AuditManifestPayload struct {
	Manifest string `json:"manifest"` // Manifest is the content to be audited.
}

// decode detects the correct decoder for use on an HTTP request and
// marshals into a given interface.
func (payload *AuditManifestPayload) Decode(r *http.Request) error {
	// Check if the content type is plain text, read it as such.
	contentType := render.GetRequestContentType(r)
	switch contentType {
	case render.ContentTypePlainText:
		// Read the request body.
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close() // Ensure the body is closed after reading.
		if err != nil {
			return errors.Wrapf(err, "failed to read the request body")
		}
		// Set the read content as the manifest payload.
		payload.Manifest = string(body)
	case render.ContentTypeJSON:
		// For non-plain text, decode the JSON body into the payload.
		if err := render.DecodeJSON(r.Body, payload); err != nil {
			return err
		}
	default:
		return errors.New("unsupported media type")
	}

	return nil
}

// Payload represents the structure for audit request data. It includes the
// manifest which is typically a string containing declarative configuration
// data that needs to be audited.
type AuditPayload struct {
	core.Locator `json:",inline"`
}

// decode detects the correct decoder for use on an HTTP request and
// marshals into a given interface.
func (payload *AuditPayload) Decode(r *http.Request) error {
	// Check if the content type is plain text, read it as such.
	contentType := render.GetRequestContentType(r)
	switch contentType {
	case render.ContentTypeJSON:
		// For non-plain text, decode the JSON body into the payload.
		if err := render.DecodeJSON(r.Body, payload); err != nil {
			return err
		}
	default:
		return errors.New("unsupported media type")
	}

	return nil
}

// AuditData represents the aggregated data of scanner issues, including the
// original list of issues and their aggregated count based on title.
type AuditData struct {
	IssueTotal    int            `json:"issueTotal"`
	ResourceTotal int            `json:"resourceTotal"`
	BySeverity    map[string]int `json:"bySeverity"`
	IssueGroups   []*IssueGroup  `json:"issueGroups"`
}

type IssueGroup struct {
	Issue    scanner.Issue  `json:"issue"`
	Locators []core.Locator `json:"locators"`
}
