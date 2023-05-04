package esstorage

import "github.com/KusionStack/karbour/pkg/search/storage"

type SearchResponse struct {
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
	Index  string            `json:"_index"`
	ID     string            `json:"_id"`
	Score  float32           `json:"_score"`
	Source *storage.Resource `json:"_source"`
}

func (s *SearchResponse) GetResources() []*storage.Resource {
	if s == nil || s.Hits == nil {
		return nil
	}
	rt := make([]*storage.Resource, len(s.Hits.Hits))
	for i, hit := range s.Hits.Hits {
		rt[i] = hit.Source
	}
	return rt
}
