package elasticsearch

import (
	"fmt"
	"sort"
	"strings"

	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/aquasecurity/esquery"
)

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
