package esstorage

import (
	"context"

	"github.com/KusionStack/karbour/pkg/search/storage"
	"github.com/aquasecurity/esquery"
)

func (s *ESClient) Search(ctx context.Context, queries []storage.Query) (*storage.SearchResult, error) {
	boolQuery := esquery.Bool()
	for _, query := range queries {
		var values []interface{}
		for _, v := range query.Values {
			values = append(values, v)
		}
		switch query.Operator {
		// TODO: support other operators
		case storage.Equals:
			boolQuery.Must(esquery.Terms(query.Key, values...))
		}
	}
	esQuery := make(map[string]interface{})
	esQuery["query"] = boolQuery.Map()
	res, err := s.searchByQuery(ctx, esQuery)
	if err != nil {
		return nil, err
	}

	rt := &storage.SearchResult{}
	rt.Total = res.Hits.Total.Value
	hits := res.Hits.Hits
	resources := make([]*storage.Resource, len(hits))
	for i := range hits {
		resources[i] = hits[i].Source
	}
	rt.Resources = resources
	return rt, nil
}
