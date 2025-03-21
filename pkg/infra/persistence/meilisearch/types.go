// Copyright The Karpor Authors.
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

//nolint:tagliatelle
package meilisearch

import (
	"fmt"
)

const (
	// maxAggSize defines the maximum number of aggregation buckets that can be returned in an aggregation query.
	maxAggSize = 10000
	// maxHitsSize defines the maximum number of search hits to be returned in a search query response.
	maxHitsSize = 1000
)

type paginationConfig struct {
	Page     int64
	PageSize int64
}

type config struct {
	filter     interface{}
	query      string
	pagination *paginationConfig
}

type Option func(*config) error

// Pagination is a functional option to set the pagination configuration.
func Pagination(page, pageSize int) Option {
	return func(c *config) error {
		if c == nil {
			return fmt.Errorf("config can't be nil")
		}
		c.pagination = &paginationConfig{
			Page:     int64(page),
			PageSize: int64(pageSize),
		}
		return nil
	}
}

// Filter is a functional option to set the filter configuration.
func Filter(filter interface{}) Option {
	return func(c *config) error {
		if c == nil {
			return fmt.Errorf("config can't be nil")
		}
		c.filter = filter
		return nil
	}
}

// Query is a functional option to set the query string.
func Query(query string) Option {
	return func(c *config) error {
		if c == nil {
			return fmt.Errorf("config can't be nil")
		}
		c.query = query
		return nil
	}
}

var ErrNotFound = &MSError{
	StatusCode: 404,
	Message:    "Object not found",
}

// MSError is an error type which represents a single ES error
type MSError struct {
	StatusCode int
	Message    string
}

// Error() method implementation for MSError, which returns the error message in a string format.
func (e *MSError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.StatusCode, e.Message)
}

// CountResponse represents the response structure for a count operation.
type CountResponse struct {
	Count int64 `json:"count"`
}

// Hits contains the hit documents and metadata from a search operation.
type Hits struct {
	Total    *Total  `json:"total"`
	MaxScore float32 `json:"max_score"`
	Hits     []*Hit  `json:"hits"`
}

// Total provides information about the total number of documents matching the search query.
type Total struct {
	Value int `json:"value,omitempty"`
}

// Hit represents a single hit document from a search operation, containing index, ID, score, and source data.
type Hit struct {
	Index  string                 `json:"_index"`
	ID     string                 `json:"_id"`
	Score  float32                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

// AggResults is assumed to be a struct that holds aggregation results.
type AggResults struct {
	Buckets []Bucket
	Total   int
}

// Bucket is assumed to be a struct that holds individual bucket data.
type Bucket struct {
	Keys  []string
	Count int
}

type SearchRequest struct {
	Query                 string      `json:"query"`
	Facets                []string    `json:"facets"`
	Page                  int         `json:"page"`
	Limit                 int64       `json:"limit"`
	Offset                int64       `json:"offset"`
	Sort                  []string    `json:"sort"`
	Filter                interface{} `json:"filters"`
	FacetFilters          []string    `json:"facetFilters"`
	AttributesToRetrieve  []string    `json:"attributesToRetrieve"`
	AttributesToHighlight []string    `json:"attributesToHighlight"`
	AttributesToCrop      []string    `json:"attributesToCrop"`
}

type SearchResponse struct {
	Hits               []interface{} `json:"hits"`
	EstimatedTotalHits int64         `json:"estimatedTotalHits,omitempty"`
	Offset             int64         `json:"offset,omitempty"`
	Limit              int64         `json:"limit,omitempty"`
	ProcessingTimeMs   int64         `json:"processingTimeMs"`
	Query              string        `json:"query"`
	FacetDistribution  interface{}   `json:"facetDistribution,omitempty"`
	TotalHits          int64         `json:"totalHits,omitempty"`
	HitsPerPage        int64         `json:"hitsPerPage,omitempty"`
	Page               int64         `json:"page,omitempty"`
	TotalPages         int64         `json:"totalPages,omitempty"`
	FacetStats         interface{}   `json:"facetStats,omitempty"`
	IndexUID           string        `json:"indexUid,omitempty"`
}
