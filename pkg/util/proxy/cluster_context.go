package proxy

import (
	"context"
)

type clusterKey int

const (
	// clusterKey is the context key for the request namespace.
	clusterContextKey clusterKey = iota
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
