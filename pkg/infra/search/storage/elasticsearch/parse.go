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

package elasticsearch

import (
	"fmt"
	"sort"
	"strings"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/elliotxx/esquery"
)

// Parse takes a query string and returns a slice of storage.Query and an error if any.
func Parse(queryString string) ([]storage.Query, error) {
	parts := splitTerms(queryString)
	sort.StringSlice(parts).Sort()
	var queries []storage.Query
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
			queries = append(queries, storage.Query{Key: lhs, Operator: storage.Equals, Values: []string{rhs}})
		default:
			return nil, fmt.Errorf("invalid query string: '%s'; can't understand '%s'", queryString, part)
		}
	}
	return queries, nil
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
func ParseQueries(queries []storage.Query) (map[string]interface{}, error) {
	boolQuery := esquery.Bool()
	for _, query := range queries {
		switch query.Operator {
		case storage.Equals:
			boolQuery.Must(esquery.Term(query.Key, query.Values[0]))
		default:
			return nil, fmt.Errorf("invalid query operator %s", query.Operator)
		}
	}
	esQuery := make(map[string]interface{})
	esQuery["query"] = boolQuery.Map()
	return esQuery, nil
}
