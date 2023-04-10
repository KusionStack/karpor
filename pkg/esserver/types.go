package esserver

const (
	GroupPath   = "group"
	VersionPath = "version"
	KindPath    = "kind"
)

type TermOptions struct {
	Must []string `yaml:"must"`
}

type QueryOptions struct {
	Version *TermOptions `yaml:"version"`
	Group   *TermOptions `yaml:"group"`
	Kind    *TermOptions `yaml:"kind"`
}

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
	Group           string                 `json:"group"`
	Version         string                 `json:"version"`
	Kind            string                 `json:"kind"`
	Resource        string                 `json:"resource"`
	ResourceVersion string                 `json:"resource_version"`
	Name            string                 `json:"name"`
	Namespace       string                 `json:"namespace"`
	Object          map[string]interface{} `json:"object"`
}
