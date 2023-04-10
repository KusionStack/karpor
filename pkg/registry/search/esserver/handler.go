package esserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticServer struct {
	client *elasticsearch.Client
}

func NewElasticServerOrDie(cfg elasticsearch.Config) *ElasticServer {
	cl, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return &ElasticServer{
		client: cl,
	}
}

func (es *ElasticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ConvertBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var searchFuns []func(searchRequest *esapi.SearchRequest)
	searchFuns = append(searchFuns, es.client.Search.WithContext(r.Context()))
	searchFuns = append(searchFuns, es.client.Search.WithIndex("clusterpedia-v1-common"))
	searchFuns = append(searchFuns, es.client.Search.WithBody(body))
	resp, err := es.client.Search(searchFuns...)
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.IsError() {
		http.Error(w, resp.String(), resp.StatusCode)
		return
	}

	var sr SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		http.Error(w, resp.String(), resp.StatusCode)
		return
	}

	b, err := json.Marshal(sr.GetResources())
	if err != nil {
		http.Error(w, resp.String(), resp.StatusCode)
		return
	}
	// buf := &bytes.Buffer{}
	// _, err = buf.ReadFrom(resp.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func ConvertBody(body io.Reader) (io.Reader, error) {
	var opts QueryOptions
	if err := json.NewDecoder(body).Decode(&opts); err != nil {
		return nil, err
	}

	builder := NewQueryBuilder()
	if opts.Version != nil {
		builder.addExpression(NewTerms(VersionPath, []string{strings.ToLower(opts.Version.Must[0])}))
	}
	if opts.Group != nil {
		builder.addExpression(NewTerms(GroupPath, []string{strings.ToLower(opts.Group.Must[0])}))
	}
	if opts.Kind != nil {
		builder.addExpression(NewTerms(KindPath, []string{strings.ToLower(opts.Kind.Must[0])}))
	}
	query := builder.build()
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %v", err)
	}
	return buf, nil
}

func (r *SearchResponse) GetResources() []*Resource {
	hits := r.Hits.Hits
	resources := make([]*Resource, len(hits))
	for i := range hits {
		resources[i] = hits[i].Source
	}
	return resources
}
