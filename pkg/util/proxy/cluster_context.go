package proxy

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

type clusterKey int

const (
	// clusterKey is the context key for the request namespace.
	clusterContextKey clusterKey = iota
	ClusterProxyURL              = "/apis/cluster.karbour.com/v1beta1/clusters/%s/proxy/"
)

// WithCluster returns a context that describes the nested cluster context
func WithCluster(parent context.Context, cluster string) context.Context {
	return context.WithValue(parent, clusterContextKey, cluster)
}

// ClusterFrom returns the value of the cluster key on the ctx
func ClusterFrom(ctx context.Context) (string, bool) {
	cluster, ok := ctx.Value(clusterContextKey).(string)
	if !ok {
		return "", false
	}
	return cluster, true
}

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
