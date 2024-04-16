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

package resourcegrouprule

import (
	"net/http"
	"time"

	"github.com/KusionStack/karbour/pkg/core/entity"
	"github.com/KusionStack/karbour/pkg/core/handler"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

// Ensure that ResourceGroupRulePayload implements the handler.Payload
// interface.
var _ handler.Payload = &ResourceGroupRulePayload{}

type ResourceGroupRulePayload struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Fields      []string `json:"fields"`
}

// decode detects the correct decoder for use on an HTTP request and
// marshals into a given interface.
func (p *ResourceGroupRulePayload) Decode(r *http.Request) error {
	// Check if the content type is plain text, read it as such.
	contentType := render.GetRequestContentType(r)
	switch contentType {
	case render.ContentTypeJSON:
		// For non-plain text, decode the JSON body into the payload.
		if err := render.DecodeJSON(r.Body, p); err != nil {
			return err
		}
	default:
		return errors.New("unsupported media type")
	}

	return nil
}

func (p *ResourceGroupRulePayload) ToEntity() *entity.ResourceGroupRule {
	return &entity.ResourceGroupRule{
		Name:        p.Name,
		Description: p.Description,
		Fields:      p.Fields,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
