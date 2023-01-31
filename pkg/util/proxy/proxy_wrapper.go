package proxy

import (
	"fmt"
	"net/http"
	"path"
)

const (
	ClusterProxyURL = "/apis/cluster.karbour.com/v1beta1/clusterextensions/%s/proxy/"
)

func WithProxyByCluster(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cluster, ok := ClusterFrom(req.Context())
		if ok && cluster != "" {
			url := fmt.Sprintf(ClusterProxyURL, cluster)
			req.URL.Path = path.Join(url, req.URL.Path)
			req.RequestURI = path.Join(url, req.RequestURI)
		}
		handler.ServeHTTP(w, req)
	})
}
