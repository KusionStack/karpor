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

//nolint:tagliatelle
package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Client represents an Elasticsearch client that can perform various operations on the Elasticsearch cluster.
type Client struct {
	client      *elasticsearch.Client
	typedClient *elasticsearch.TypedClient
}

// NewClient creates a new Elasticsearch client instance
func NewClient(config elasticsearch.Config) (*Client, error) {
	cl, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	typed, err := elasticsearch.NewTypedClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{client: cl, typedClient: typed}, nil
}

// SaveDocument saves a new document
func (cl *Client) SaveDocument(
	ctx context.Context,
	indexName string,
	documentID string,
	body io.Reader,
) error {
	resp, err := cl.client.Index(
		indexName,
		body,
		cl.client.Index.WithDocumentID(documentID),
		cl.client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return &ESError{
			StatusCode: resp.StatusCode,
			Message:    resp.String(),
		}
	}
	return nil
}

// GetDocument gets a document with the specified ID
func (cl *Client) GetDocument(
	ctx context.Context,
	indexName string,
	documentID string,
) (map[string]interface{}, error) {
	resp, err := cl.client.Get(indexName, documentID, cl.client.Get.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, &ESError{
			StatusCode: resp.StatusCode,
			Message:    resp.String(),
		}
	}
	getResp := &struct {
		Index  string                 `json:"_index"`
		ID     string                 `json:"_id"`
		Found  bool                   `json:"found"`
		Source map[string]interface{} `json:"_source"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(getResp); err != nil {
		return nil, err
	}
	if !getResp.Found {
		return nil, ErrNotFound
	}
	return getResp.Source, nil
}

// DeleteDocument deletes a document with the specified ID
func (cl *Client) DeleteDocument(ctx context.Context, indexName string, documentID string) error {
	if _, err := cl.GetDocument(ctx, indexName, documentID); err != nil {
		return err
	}
	resp, err := cl.client.Delete(indexName, documentID, cl.client.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return &ESError{
			StatusCode: resp.StatusCode,
			Message:    resp.String(),
		}
	}
	return nil
}

// DeleteDocumentByQuery deletes documents from the specified index based on the
// provided query in the body.
func (cl *Client) DeleteDocumentByQuery(
	ctx context.Context,
	indexName string,
	body io.Reader,
) error {
	resp, err := cl.client.DeleteByQuery(
		[]string{indexName},
		body,
		cl.client.DeleteByQuery.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return &ESError{
			StatusCode: resp.StatusCode,
			Message:    resp.String(),
		}
	}
	return nil
}

// SearchDocument performs a search query in the specified index
func (cl *Client) SearchDocument(
	ctx context.Context,
	indexName string,
	body io.Reader,
	options ...Option,
) (*SearchResponse, error) {
	cfg := &config{}
	for _, option := range options {
		if err := option(cfg); err != nil {
			return nil, err
		}
	}
	opts := []func(*esapi.SearchRequest){
		cl.client.Search.WithContext(ctx),
		cl.client.Search.WithIndex(indexName),
		cl.client.Search.WithBody(body),
	}
	if cfg.pagination != nil {
		from := (cfg.pagination.Page - 1) * cfg.pagination.PageSize
		opts = append(
			opts,
			cl.client.Search.WithSize(cfg.pagination.PageSize),
			cl.client.Search.WithFrom(from),
		)
	}

	resp, err := cl.client.Search(opts...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return nil, &ESError{
			StatusCode: resp.StatusCode,
			Message:    resp.String(),
		}
	}

	sr := &SearchResponse{}
	if err := json.NewDecoder(resp.Body).Decode(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

// CreateIndex creates a new index with the specified settings and mappings
func (cl *Client) CreateIndex(ctx context.Context, index string, body io.Reader) error {
	resp, err := cl.client.Indices.Create(
		index,
		cl.client.Indices.Create.WithBody(body),
		cl.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		if strings.Contains(resp.String(), "resource_already_exists_exception") {
			return nil
		}
		return &ESError{
			StatusCode: resp.StatusCode,
			Message:    resp.String(),
		}
	}
	return nil
}

// IsIndexExists Check if an index exists in Elasticsearch
func (cl *Client) IsIndexExists(ctx context.Context, index string) (bool, error) {
	resp, err := cl.client.Indices.Exists(
		[]string{index},
		cl.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return false, &ESError{
			StatusCode: resp.StatusCode,
			Message:    resp.String(),
		}
	}
	// Decide if the index exists based on the response status code
	if resp.StatusCode == http.StatusOK {
		return true, nil // Index exists
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil // Index does not exist
	} else {
		// If it's any other status code, return an error
		return false, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}
}

// AggregationByTerms performs an aggregation query based on the provided fields.
func (cl *Client) AggregationByTerms(ctx context.Context, index string, name string, fields []string) (*AggResults, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields provided for aggregation")
	}
	if len(fields) == 1 {
		// Perform single-term aggregation if only one field is provided.
		return cl.termsAggSearch(ctx, index, name, fields[0])
	}
	// Perform multi-term aggregation if multiple fields are provided.
	return cl.multiTermsAggSearch(ctx, index, name, fields)
}

// multiTermsAggSearch executes a multi-term aggregation query on specified fields.
func (cl *Client) multiTermsAggSearch(ctx context.Context, index string, name string, fields []string) (*AggResults, error) {
	// Construct the terms for multi-term aggregation based on the fields.
	terms := make([]types.MultiTermLookup, len(fields))
	for i := range fields {
		terms[i] = types.MultiTermLookup{Field: fields[i]}
	}

	// Execute the search request with the constructed multi-term aggregation.
	resp, err := cl.typedClient.
		Search().
		Index(index).
		Request(&search.Request{
			// Set the number of search hits to return to 0 as we only need aggregation data.
			Size: some.Int(0),
			Aggregations: map[string]types.Aggregations{
				name: {
					MultiTerms: &types.MultiTermsAggregation{
						Terms: terms,
						// maxAggSize should be predefined to limit the size of the aggregation.
						Size: some.Int(maxAggSize),
					},
				},
			},
		}).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// Extract the buckets from the response and construct the AggResults.
	buckets := resp.Aggregations[name].(*types.MultiTermsAggregate).Buckets.([]types.MultiTermsBucket)
	bs := make([]Bucket, len(buckets))
	for i, b := range buckets {
		keys := make([]string, len(b.Key))
		for j, k := range b.Key {
			keys[j] = fmt.Sprintf("%v", k)
		}
		bs[i] = Bucket{
			Keys:  keys,
			Count: int(b.DocCount),
		}
	}
	return &AggResults{
		Buckets: bs,
		Total:   len(bs),
	}, nil
}

// termsAggSearch executes a single-term aggregation query on the specified field.
func (cl *Client) termsAggSearch(ctx context.Context, index string, name string, field string) (*AggResults, error) {
	// Execute the search request with the single-term aggregation.
	resp, err := cl.typedClient.
		Search().
		Index(index).
		Request(&search.Request{
			// Set the number of search hits to return to 0 as we only need aggregation data.
			Size: some.Int(0),
			Aggregations: map[string]types.Aggregations{
				name: {
					Terms: &types.TermsAggregation{
						Field: some.String(field),
						// maxAggSize should be predefined to limit the size of the aggregation.
						Size: some.Int(maxAggSize),
					},
				},
			},
		}).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// Extract the buckets from the response and construct the AggResults.
	buckets := resp.Aggregations[name].(*types.StringTermsAggregate).Buckets.([]types.StringTermsBucket)
	bs := make([]Bucket, len(buckets))
	for i, b := range buckets {
		bs[i] = Bucket{
			Keys:  []string{fmt.Sprintf("%v", b.Key)},
			Count: int(b.DocCount),
		}
	}
	return &AggResults{
		Buckets: bs,
		Total:   len(bs),
	}, nil
}
