// Copyright The Karbour Authors.
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

package proxy

import (
	"context"
	"fmt"
	"net/http"
	"path"

	filtersutil "github.com/KusionStack/karbour/pkg/util/filters"
)

type clusterKey int

const (
	// clusterKey is the context key for the request namespace.
	clusterContextKey   clusterKey = iota
	ClusterProxyURL                = "/apis/cluster.karbour.com/v1beta1/clusters/%s/proxy/"
	clusterNamespaceKey            = "namespace"
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

func NamespaceFrom(ctx context.Context) (string, bool) {
	namespace, ok := ctx.Value(filtersutil.CtxTyp(clusterNamespaceKey)).(string)
	if !ok {
		return "", false
	}
	return namespace, true
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
