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
	// 对应 ES 的 max_result_window
	Pagination: &meilisearch.Pagination{
		MaxTotalHits: 1_000_000, // 控制最大返回结果数
	},

	// 对应 ES analysis.normalizer
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

	// 对应 ES mappings.properties 字段配置
	FilterableAttributes: []string{
		"id",
		"cluster",
		"apiVersion",
		"kind", // 启用过滤和精确匹配
		"namespace",
		"name",
		"labels", // 扁平化字段过滤
		"annotations",
		"ownerReferences",
		"content",
		"resourceVersion",
		"syncAt",
		"deleted",
	},

	SearchableAttributes: []string{
		"content", // 全文搜索字段
	},

	// 对应 _source.excludes 排除字段
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
		// 排除 "custom" 字段
	},
}
var defaultResourceGroupRuleIndexSetting = &meilisearch.Settings{
	// 对应 Elasticsearch 的 max_result_window
	Pagination: &meilisearch.Pagination{
		MaxTotalHits: 1000000, // 控制最大返回结果数
	}, // 控制最大返回结果数

	// 对应 mappings.properties 字段属性
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

	// 对应 mappings._source.excludes
	// Meilisearch 无完全等效功能，但可通过以下方式近似实现
	DisplayedAttributes: []string{
		"id",
		"name",
		"description",
		"fields",
		"createdAt",
		"updatedAt",
		"deletedAt",
	}, // 显式指定返回字段

}
