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
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
)

// Parse takes a query string and returns a slice of storage.Query and an error if any.
func Parse(queryString string) ([]interface{}, error) {
	parts := splitTerms(queryString)
	sort.StringSlice(parts).Sort()
	var filters []interface{}
	for _, part := range parts {
		if part == "" {
			continue
		}
		lhs, op, rhs, ok := splitTerm(part)
		if !ok {
			return nil, fmt.Errorf("invalid query string: '%s'; can't understand '%s'", queryString, part)
		}

		switch op {
		case storage.Equals:
			if len(rhs) == 0 {
				return nil, fmt.Errorf("invalid query string: '%s'; can't understand '%s'", queryString, part)
			}
			strs := strings.Split(rhs, ",")
			if len(strs) == 0 {
				return nil, fmt.Errorf("invalid query string: '%s'; can't understand '%s'", queryString, part)
			}
			if len(strs) == 1 {
				filters = append(filters, fmt.Sprintf("%s = %s", lhs, strs[0]))
			} else {
				filters = append(filters, fmt.Sprintf("%s IN [%s]", lhs, strings.Join(strs, ",")))
			}
		default:
			return nil, fmt.Errorf("invalid query string: '%s'; can't understand '%s'", queryString, part)
		}
	}
	return filters, nil
}

// splitTerms returns the comma-separated terms contained in the given fieldSelector.
// Backslash-escaped commas are treated as data instead of delimiters, and are included in the returned terms, with the leading backslash preserved.
func splitTerms(s string) []string {
	if len(s) == 0 {
		return nil
	}

	terms := make([]string, 0, 1)
	startIndex := 0
	inSlash := false
	for i, c := range s {
		switch {
		case inSlash:
			inSlash = false
		case c == '\\':
			inSlash = true
		case c == ',':
			terms = append(terms, s[startIndex:i])
			startIndex = i + 1
		}
	}
	terms = append(terms, s[startIndex:])
	return terms
}

// termOperators holds the recognized operators supported in fieldSelectors.
var termOperators = []string{storage.Equals}

// splitTerm returns the lhs, operator, and rhs parsed from the given term, along with an indicator of whether the parse was successful.
func splitTerm(term string) (lhs, op, rhs string, ok bool) {
	for i := range term {
		remaining := term[i:]
		for _, op := range termOperators {
			if strings.HasPrefix(remaining, op) {
				return term[0:i], op, term[i+len(op):], true
			}
		}
	}
	return "", "", "", false
}

// ParseQueries takes a slice of storage.Query and returns a map of interface{} representing the parsed queries and an error if any.
func ParseQueries(queries []storage.Query) ([]interface{}, error) {
	var filter []interface{}
	for _, query := range queries {
		switch query.Operator {
		case storage.Equals:
			if len(query.Values) == 0 {
				return nil, fmt.Errorf("invalid query: %s", query.Key)
			}
			if len(query.Values) == 1 {
				filter = append(filter, fmt.Sprintf("%s = %s", query.Key, query.Values[0]))
			} else {
				filter = append(filter, fmt.Sprintf("%s IN [%s]", query.Key, strings.Join(query.Values, ",")))
			}
		default:
			return nil, fmt.Errorf("invalid query operator %s", query.Operator)
		}
	}
	return filter, nil
}

// ConvertToFilter 将 map 转换为 Meilisearch 等值过滤器
// 输入示例:
//
//	{
//	  "category": "electronics",
//	  "tags": ["a", "b"],
//	  "author": {"name": "Alice"}
//	}
//
// 输出: ["category = 'electronics'", "tags IN ['a','b']", "author.name = 'Alice'"]
func ConvertToFilter(query map[string]any) ([]any, error) {
	return processMap(query, "")
}

// 递归处理嵌套 map
func processMap(data map[string]any, parentKey string) ([]any, error) {
	var filters []any

	for key, value := range data {
		fullKey := key
		if parentKey != "" {
			fullKey = parentKey + "." + key
		}

		switch v := value.(type) {
		case map[string]any:
			// 处理嵌套对象
			nestedFilters, err := processMap(v, fullKey)
			if err != nil {
				return nil, err
			}
			filters = append(filters, nestedFilters)

		case []any:
			// 处理数组 (IN 查询)
			inClause, err := buildInClause(fullKey, v)
			if err != nil {
				return nil, err
			}
			filters = append(filters, inClause)

		default:
			// 处理基本类型
			eqClause, err := buildEqualClause(fullKey, v)
			if err != nil {
				return nil, err
			}
			filters = append(filters, eqClause)
		}
	}

	return filters, nil
}

// 构建等值条件
func buildEqualClause(key string, value any) (string, error) {
	valueStr, err := convertValue(value)
	if err != nil {
		return "", fmt.Errorf("key %s: %v", key, err)
	}
	return fmt.Sprintf("%s = %s", key, valueStr), nil
}

// 构建 IN 条件
func buildInClause(key string, values []any) (string, error) {
	var elements []string
	for _, v := range values {
		element, err := convertValue(v)
		if err != nil {
			return "", fmt.Errorf("key %s: %v", key, err)
		}
		elements = append(elements, element)
	}
	return fmt.Sprintf("%s IN [%s]", key, strings.Join(elements, ",")), nil
}

// 值类型转换
func convertValue(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "\\'") + "'", nil
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("unsupported type %T", value)
	}
}
