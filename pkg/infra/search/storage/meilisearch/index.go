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

import "github.com/meilisearch/meilisearch-go"

const (
	defaultResourceIndexName          = "resources"
	defaultResourceGroupRuleIndexName = "resource_group_rules"
)

var defaultResourceIndexSetting = &meilisearch.Settings{
	Pagination: &meilisearch.Pagination{
		MaxTotalHits: 1_000_000,
	},

	SortableAttributes: []string{
		"creationTimestamp",
		"deletionTimestamp",
	},
	Faceting: &meilisearch.Faceting{
		MaxValuesPerFacet: 1000,
		SortFacetValuesBy: map[string]meilisearch.SortFacetType{
			"cluster":    meilisearch.SortFacetTypeAlpha,
			"apiVersion": meilisearch.SortFacetTypeAlpha,
			"kind":       meilisearch.SortFacetTypeAlpha,
			"namespace":  meilisearch.SortFacetTypeAlpha,
		},
	},

	FilterableAttributes: []string{
		"id",
		"cluster",
		"apiVersion",
		"kind",
		"namespace",
		"name",
		"labels",
		"annotations",
		"ownerReferences",
		"content",
		"resourceVersion",
		"syncAt",
		"deleted",
	},

	SearchableAttributes: []string{
		"content",
	},

	DisplayedAttributes: []string{
		"cluster",
		"apiVersion",
		"kind",
		"namespace",
		"name",
		"creationTimestamp",
		"deletionTimestamp",
		"labels",
		"annotations",
		"ownerReferences",
		"resourceVersion",
		"content",
	},
}

var defaultResourceGroupRuleIndexSetting = &meilisearch.Settings{
	Pagination: &meilisearch.Pagination{
		MaxTotalHits: 1000000,
	},

	FilterableAttributes: []string{
		"id",
		"name",
		"description",
		"fields",
		"deleted",
	},
	SortableAttributes: []string{
		"createdAt",
		"updatedAt",
		"deletedAt",
	},

	DisplayedAttributes: []string{
		"id",
		"name",
		"description",
		"fields",
		"createdAt",
		"updatedAt",
		"deletedAt",
	},
}
