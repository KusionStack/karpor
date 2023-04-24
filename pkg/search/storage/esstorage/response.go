package esstorage

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
	Index  string    `json:"_index"`
	ID     string    `json:"_id"`
	Score  float32   `json:"_score"`
	Source *Resource `json:"_source"`
}

type Resource struct {
	Cluster    string                 `json:"cluster"`
	Namespace  string                 `json:"namespace"`
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Name       string                 `json:"name"`
	Object     map[string]interface{} `json:"object"`
}

func (r *SearchResponse) GetResources() []*Resource {
	hits := r.Hits.Hits
	resources := make([]*Resource, len(hits))
	for i := range hits {
		resources[i] = hits[i].Source
	}
	return resources
}
