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

package meilisearch

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

// Client represents an Elasticsearch client that can perform various operations on the Elasticsearch cluster.
type Client struct {
	client meilisearch.ServiceManager
}

// NewClient creates a new Elasticsearch client instance
func NewClient(address string, key string) (*Client, error) {
	cl := meilisearch.New(address, meilisearch.WithAPIKey(key))
	return &Client{client: cl}, nil
}

// SaveDocument saves a new document
func (cl *Client) SaveDocument(
	ctx context.Context,
	indexName string,
	obj interface{},
) error {
	taskInfo, err := cl.client.Index(indexName).UpdateDocumentsWithContext(ctx, obj)
	if err != nil {
		return err
	}
	return cl.WaitForTask(ctx, taskInfo)
}

// GetDocument gets a document with the specified ID
func (cl *Client) GetDocument(
	ctx context.Context,
	indexName string,
	documentID string,
) (map[string]interface{}, error) {
	getResp := map[string]interface{}{}
	err := cl.client.Index(indexName).GetDocumentWithContext(ctx, documentID, nil, &getResp)
	if err != nil {
		return nil, err
	}

	return getResp, nil
}

// DeleteDocument deletes a document with the specified ID
func (cl *Client) DeleteDocument(ctx context.Context, indexName, documentID string) error {
	if _, err := cl.GetDocument(ctx, indexName, documentID); err != nil {
		return err
	}
	resp, err := cl.client.Index(indexName).DeleteDocumentWithContext(ctx, documentID)
	if err != nil {
		return err
	}
	return cl.WaitForTask(ctx, resp)
}

// DeleteDocumentByQuery deletes documents from the specified index based on the
// provided query in the body.
func (cl *Client) DeleteDocumentByQuery(
	ctx context.Context,
	indexName string,
	filter interface{},
) error {
	task, err := cl.client.Index(indexName).DeleteDocumentsByFilter(filter)
	if err != nil {
		return err
	}
	return cl.WaitForTask(ctx, task)
}

// SearchDocument performs a search query in the specified index
func (cl *Client) SearchDocument(ctx context.Context, indexName string, searchRequest *SearchRequest) (*SearchResponse, error) {
	req := &meilisearch.SearchRequest{
		Query:                 searchRequest.Query,
		Facets:                searchRequest.Facets,
		Limit:                 searchRequest.Limit,
		Offset:                searchRequest.Offset,
		Sort:                  searchRequest.Sort,
		Filter:                searchRequest.Filter,
		AttributesToRetrieve:  searchRequest.AttributesToRetrieve,
		AttributesToHighlight: searchRequest.AttributesToHighlight,
		AttributesToCrop:      searchRequest.AttributesToCrop,
	}
	resp, err := cl.client.Index(indexName).SearchWithContext(ctx, "", req)
	if err != nil {
		return nil, err
	}
	return &SearchResponse{
		Hits:               resp.Hits,
		TotalHits:          resp.TotalHits,
		Offset:             resp.Offset,
		Limit:              resp.Limit,
		ProcessingTimeMs:   resp.ProcessingTimeMs,
		Query:              resp.Query,
		FacetDistribution:  resp.FacetDistribution,
		IndexUID:           resp.IndexUID,
		FacetStats:         resp.FacetStats,
		TotalPages:         resp.TotalPages,
		HitsPerPage:        resp.HitsPerPage,
		Page:               resp.Page,
		EstimatedTotalHits: resp.EstimatedTotalHits,
	}, nil
}

// Count performs a count query in the specified index.
func (cl *Client) Count(
	ctx context.Context,
	indexName string,
) (*CountResponse, error) {
	resp, err := cl.client.Index(indexName).GetStatsWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &CountResponse{
		Count: resp.NumberOfDocuments,
	}, nil
}

// CreateIndex creates a new index with the specified settings and mappings,PrimaryKey is id by default
func (cl *Client) CreateIndex(ctx context.Context, index string, settings *meilisearch.Settings) error {
	exist, err := cl.IsIndexExists(ctx, index)
	if err != nil {
		return err
	}
	if !exist {
		resp, err := cl.client.CreateIndexWithContext(ctx, &meilisearch.IndexConfig{
			Uid: index,
		})
		if err != nil {
			return err
		}
		err = cl.WaitForTask(ctx, resp)
		if err != nil {
			return err
		}
	}
	if settings != nil {
		task, err := cl.client.Index(index).UpdateSettingsWithContext(ctx, settings)
		if err != nil {
			return err
		}
		return cl.WaitForTask(ctx, task)
	}
	return nil
}

// IsIndexExists Check if an index exists in Elasticsearch
func (cl *Client) IsIndexExists(ctx context.Context, index string) (bool, error) {
	_, err := cl.client.GetIndexWithContext(ctx, index)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AggregateDocumentByTerms performs an aggregation query based on the provided fields.
func (cl *Client) AggregateDocumentByTerms(ctx context.Context, index string, fields []string) (*AggResults, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields provided for aggregation")
	}
	if len(fields) > 1 {
		return nil, fmt.Errorf("only one field is supported for aggregation for meilisearch")
	}
	// Execute the search request with the single-term aggregation.
	resp, err := cl.client.Index(index).SearchWithContext(ctx, "", &meilisearch.SearchRequest{
		Facets: fields,
		Limit:  1, // no hits needed
	})
	if err != nil {
		return nil, err
	}
	aggRes := resp.FacetDistribution.(map[string]interface{})
	filedAgg := aggRes[fields[0]].(map[string]interface{})
	results := &AggResults{
		Total: 1,
	}

	for fieldValue, count := range filedAgg {
		var bucket Bucket
		bucket.Keys = []string{fieldValue}
		bucket.Count = int(count.(float64))
		results.Total++
		results.Buckets = append(results.Buckets, bucket)
	}
	return results, nil
}

func (cl *Client) UpdateIndexFacets(ctx context.Context, index string, fields []string) error {
	facets, err := cl.client.Index(index).GetFacetingWithContext(ctx)
	if err != nil {
		return err
	}
	if facets == nil {
		facets = &meilisearch.Faceting{
			MaxValuesPerFacet: maxAggSize,
		}
	}
	if facets.SortFacetValuesBy == nil {
		facets.SortFacetValuesBy = make(map[string]meilisearch.SortFacetType)
	}
	for _, field := range fields {
		facets.SortFacetValuesBy[field] = meilisearch.SortFacetTypeAlpha
	}
	taskInfo, err := cl.client.Index(index).UpdateFacetingWithContext(ctx, facets)
	if err != nil {
		return err
	}
	return cl.WaitForTask(ctx, taskInfo)
}

// SearchLiveness would ping ElasticSearch client
func (cl *Client) SearchLiveness(_ context.Context) error {
	_, err := cl.client.Health()
	return err
}

func (cl *Client) WaitForTask(ctx context.Context, taskInfo *meilisearch.TaskInfo) error {
	if taskInfo.Status == meilisearch.TaskStatusSucceeded {
		return nil
	}
	task, err := cl.client.WaitForTaskWithContext(ctx, taskInfo.TaskUID, 0)
	if err != nil {
		return err
	}
	if task.Status != meilisearch.TaskStatusSucceeded {
		return fmt.Errorf("task status is not Succeeded: %s %s,%v", task.Status, task.Error, task.Details)
	}
	return nil
}
