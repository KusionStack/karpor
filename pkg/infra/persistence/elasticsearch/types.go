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

package elasticsearch

import (
	"context"
	"fmt"
	"io"
)

// Client defines the interface for our Elasticsearch operations.
type Client interface {
	Get(ctx context.Context, index string, documentID string) (any, error)
	Save(ctx context.Context, index string, documentID string, body io.Reader) error
	Delete(ctx context.Context, index string, documentID string) error
	Search(ctx context.Context, index string, body io.Reader, options ...Option) (*SearchResp, error)
	CreateIndex(ctx context.Context, index string, body io.Reader) error
	ExistIndex(ctx context.Context, index string) (bool, error)
}

type paginationConfig struct {
	Page     int
	PageSize int
}

type config struct {
	pagination *paginationConfig
}

type Option func(*config) error

func Pagination(page, pageSize int) Option {
	return func(c *config) error {
		if c == nil {
			return fmt.Errorf("config can't be nil")
		}
		c.pagination = &paginationConfig{
			Page:     page,
			PageSize: pageSize,
		}
		return nil
	}
}

var ESErrorNotFound = &ESError{
	StatusCode: 404,
	Message:    "Object not found",
}

// ESError is an error type which represents a single ES error
type ESError struct {
	StatusCode int
	Message    string
}

func (e *ESError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.StatusCode, e.Message)
}

type SearchResp struct {
	ScrollID string `json:"_scroll_id"`
	Took     int    `json:"took"`
	TimeOut  bool   `json:"time_out"`
	Hits     *Hits  `json:"hits"`
}

type Hits struct {
	Total    *Total  `json:"total"`
	MaxScore float32 `json:"max_score"`
	Hits     []*Hit  `json:"hits"`
}

type Total struct {
	Value    int    `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type Hit struct {
	Index  string  `json:"_index"`
	ID     string  `json:"_id"`
	Score  float32 `json:"_score"`
	Source any     `json:"_source"`
}
