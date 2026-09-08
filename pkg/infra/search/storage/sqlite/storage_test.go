// Copyright The Karpor Authors.
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

package sqlite

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

func TestStorage_InitializesAndReopens(t *testing.T) {
	dsn := testDSN(t)

	s, err := NewStorage(&Config{DataSourceName: dsn})
	require.NoError(t, err)

	count, err := s.CountResourceGroupRules(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, count)

	rule, err := s.GetResourceGroupRule(context.Background(), defaultResourceGroupRuleName)
	require.NoError(t, err)
	require.Equal(t, []string{defaultResourceGroupRuleName}, rule.Fields)
	require.NoError(t, s.Close())

	reopened, err := NewStorage(&Config{DataSourceName: dsn})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reopened.Close()) })

	count, err = reopened.CountResourceGroupRules(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestStorage_SaveGetUpdateResource(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	pod := testResource("uid-1", "v1", "Pod", "default", "api", map[string]string{"app": "api"}, map[string]string{"owner": "team-a"})
	pod.SetResourceVersion("1")
	require.NoError(t, s.SaveResource(ctx, "cluster-a", pod))

	got := testDefaultPodLookup("api")
	require.NoError(t, s.GetResource(ctx, "cluster-a", got))
	require.Equal(t, "api", got.GetName())
	require.Equal(t, "api", got.GetLabels()["app"])
	require.Equal(t, "team-a", got.GetAnnotations()["owner"])
	require.Equal(t, "1", got.GetResourceVersion())

	pod.SetLabels(map[string]string{"app": "api-v2"})
	pod.SetResourceVersion("2")
	require.NoError(t, s.SaveResource(ctx, "cluster-a", pod))

	got = testDefaultPodLookup("api")
	require.NoError(t, s.GetResource(ctx, "cluster-a", got))
	require.Equal(t, "api-v2", got.GetLabels()["app"])
	require.Equal(t, "2", got.GetResourceVersion())

	count, err := s.CountResources(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestStorage_DeleteResourceAndDeleteAllResources(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	require.NoError(t, s.SaveResource(ctx, "cluster-a", testResource("uid-1", "v1", "Pod", "default", "api", nil, nil)))
	require.NoError(t, s.SaveResource(ctx, "cluster-a", testResource("uid-2", "v1", "Service", "default", "api", nil, nil)))
	require.NoError(t, s.SaveResource(ctx, "cluster-b", testResource("uid-3", "v1", "Pod", "default", "api", nil, nil)))

	require.NoError(t, s.DeleteResource(ctx, "cluster-a", testDefaultPodLookup("api")))

	count, err := s.CountResources(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	require.NoError(t, s.DeleteAllResources(ctx, "cluster-a"))

	count, err = s.CountResources(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	res, err := s.SearchByTerms(ctx, map[string]any{"cluster": "cluster-b"}, nil)
	require.NoError(t, err)
	require.Equal(t, 1, res.Total)
	require.Equal(t, "cluster-b", res.Resources[0].Cluster)
}

func TestStorage_SoftDeleteAndSearchFiltering(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	require.NoError(t, s.SaveResource(ctx, "cluster-a", testResource("uid-1", "v1", "Pod", "default", "active", nil, nil)))
	require.NoError(t, s.SaveResource(ctx, "cluster-a", testResource("uid-2", "v1", "Pod", "default", "deleted", nil, nil)))
	require.NoError(t, s.SoftDeleteResource(ctx, "cluster-a", testDefaultPodLookup("deleted")))

	count, err := s.CountResources(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	active, err := s.Search(ctx, "SELECT * from resources WHERE kind = 'Pod'", storage.SQLPatternType, nil)
	require.NoError(t, err)
	require.Equal(t, 1, active.Total)
	require.Equal(t, "active", active.Resources[0].Name)

	deleted, err := s.Search(ctx, "SELECT * from resources WHERE deleted = 'true'", storage.SQLPatternType, nil)
	require.NoError(t, err)
	require.Equal(t, 1, deleted.Total)
	require.True(t, deleted.Resources[0].Deleted)

	deleted, err = s.Search(ctx, "SELECT * from resources WHERE deleted = true", storage.SQLPatternType, nil)
	require.NoError(t, err)
	require.Equal(t, 1, deleted.Total)
	require.True(t, deleted.Resources[0].Deleted)
}

func TestStorage_SearchSQLSubset(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()
	seedSearchResources(t, s)

	tests := []struct {
		name      string
		query     string
		total     int
		wantNames []string
	}{
		{
			name:      "select all active resources",
			query:     "SELECT * from resources",
			total:     4,
			wantNames: []string{"api", "web", "web", "db"},
		},
		{
			name:      "and scalar fields",
			query:     "SELECT * from resources WHERE cluster = 'cluster-a' AND kind = 'Pod'",
			total:     2,
			wantNames: []string{"api", "web"},
		},
		{
			name:      "label field",
			query:     "SELECT * from resources WHERE `labels.app.kubernetes.io/name` = 'demo'",
			total:     1,
			wantNames: []string{"api"},
		},
		{
			name:      "annotation field",
			query:     "SELECT * from resources WHERE `annotations.owner` = 'team-b'",
			total:     1,
			wantNames: []string{"web"},
		},
		{
			name:      "like keyword expansion shape",
			query:     "SELECT * from resources WHERE cluster = 'cluster-a' AND (name LIKE '%ap%' OR namespace LIKE '%ops%')",
			total:     2,
			wantNames: []string{"api", "web"},
		},
		{
			name:      "in expression",
			query:     "SELECT * from resources WHERE kind IN ('Service', 'StatefulSet')",
			total:     2,
			wantNames: []string{"web", "db"},
		},
		{
			name:      "not in expression",
			query:     "SELECT * from resources WHERE kind NOT IN ('Pod', 'Service')",
			total:     1,
			wantNames: []string{"db"},
		},
		{
			name:      "contains function",
			query:     "SELECT * from resources WHERE contains(content, 'StatefulSet')",
			total:     1,
			wantNames: []string{"db"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := s.Search(ctx, tt.query, storage.SQLPatternType, nil)
			require.NoError(t, err)
			require.Equal(t, tt.total, res.Total)
			require.Equal(t, tt.wantNames, resourceNames(res.Resources))
		})
	}
}

func TestStorage_SearchPaginationAndLimit(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()
	seedSearchResources(t, s)

	res, err := s.Search(ctx, "SELECT * from resources WHERE cluster = 'cluster-a'", storage.SQLPatternType, &storage.Pagination{Page: 2, PageSize: 1})
	require.NoError(t, err)
	require.Equal(t, 3, res.Total)
	require.Len(t, res.Resources, 1)
	require.Equal(t, "web", res.Resources[0].Name)

	res, err = s.Search(ctx, "SELECT * from resources WHERE cluster = 'cluster-a' LIMIT 2", storage.SQLPatternType, nil)
	require.NoError(t, err)
	require.Equal(t, 3, res.Total)
	require.Len(t, res.Resources, 2)
}

func TestStorage_SearchByTerms(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()
	seedSearchResources(t, s)

	res, err := s.SearchByTerms(ctx, map[string]any{
		"cluster": "cluster-a",
		"kind":    "Pod",
		"labels":  map[string]string{"app": "web"},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, 1, res.Total)
	require.Equal(t, "web", res.Resources[0].Name)

	res, err = s.SearchByTerms(ctx, map[string]any{
		"annotations.owner": "team-a",
	}, &storage.Pagination{Page: 1, PageSize: 1})
	require.NoError(t, err)
	require.Equal(t, 2, res.Total)
	require.Len(t, res.Resources, 1)
}

func TestStorage_AggregateByTerms(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()
	seedSearchResources(t, s)

	byKind, err := s.AggregateByTerms(ctx, []string{"kind"})
	require.NoError(t, err)
	require.Contains(t, byKind.Buckets, storage.Bucket{Keys: []string{"Pod"}, Count: 2})
	require.Contains(t, byKind.Buckets, storage.Bucket{Keys: []string{"Service"}, Count: 1})

	byLabel, err := s.AggregateByTerms(ctx, []string{"labels.app"})
	require.NoError(t, err)
	require.ElementsMatch(t, []storage.Bucket{
		{Keys: []string{"web"}, Count: 2},
		{Keys: []string{"api"}, Count: 1},
		{Keys: []string{"db"}, Count: 1},
	}, byLabel.Buckets)

	byClusterAndKind, err := s.AggregateByTerms(ctx, []string{"cluster", "kind"})
	require.NoError(t, err)
	require.Contains(t, byClusterAndKind.Buckets, storage.Bucket{Keys: []string{"cluster-a", "Pod"}, Count: 2})
}

func TestStorage_ResourceGroupRuleCRUDAndGrouping(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()
	seedSearchResources(t, s)

	now := metav1.Now()
	rule := &entity.ResourceGroupRule{
		Name:        "cluster-env",
		Description: "group by cluster and env",
		Fields:      []string{"cluster", "labels.env"},
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	require.NoError(t, s.SaveResourceGroupRule(ctx, rule))

	got, err := s.GetResourceGroupRule(ctx, "cluster-env")
	require.NoError(t, err)
	require.NotEmpty(t, got.ID)
	require.Equal(t, rule.Name, got.Name)
	require.Equal(t, rule.Fields, got.Fields)

	rules, err := s.ListResourceGroupRules(ctx)
	require.NoError(t, err)
	require.Len(t, rules, 2)

	count, err := s.CountResourceGroupRules(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	groups, err := s.ListResourceGroupsBy(ctx, "cluster-env")
	require.NoError(t, err)
	require.Equal(t, []string{"cluster", "labels.env"}, groups.Fields)
	require.Contains(t, groups.Groups, &entity.ResourceGroup{Cluster: "cluster-a", Labels: map[string]string{"env": "prod"}})
	require.Contains(t, groups.Groups, &entity.ResourceGroup{Cluster: "cluster-a", Labels: map[string]string{"env": "dev"}})

	require.NoError(t, s.DeleteResourceGroupRule(ctx, "cluster-env"))
	_, err = s.GetResourceGroupRule(ctx, "cluster-env")
	require.ErrorIs(t, err, storage.ErrResourceGroupRuleNotFound)
}

func TestStorage_SearchRejectsUnsafeOrUnsupportedSQL(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	tests := []struct {
		name        string
		query       string
		patternType string
	}{
		{name: "delete", query: "DELETE FROM resources", patternType: storage.SQLPatternType},
		{name: "multiple statements", query: "SELECT * from resources; DELETE FROM resources", patternType: storage.SQLPatternType},
		{name: "wrong table", query: "SELECT * from other WHERE name = 'api'", patternType: storage.SQLPatternType},
		{name: "join", query: "SELECT * from resources JOIN other ON resources.name = other.name WHERE resources.name = 'api'", patternType: storage.SQLPatternType},
		{name: "unknown field", query: "SELECT * from resources WHERE unsupported = 'api'", patternType: storage.SQLPatternType},
		{name: "subquery", query: "SELECT * from resources WHERE name IN (SELECT name FROM resources)", patternType: storage.SQLPatternType},
		{name: "unsafe function", query: "SELECT * from resources WHERE random() = 1", patternType: storage.SQLPatternType},
		{name: "dsl", query: "kind=Pod", patternType: storage.DSLPatternType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Search(ctx, tt.query, tt.patternType, nil)
			require.Error(t, err)
		})
	}
}

func TestStorage_PersistsResourcesAcrossReopen(t *testing.T) {
	dsn := testDSN(t)
	ctx := context.Background()

	s, err := NewStorage(&Config{DataSourceName: dsn})
	require.NoError(t, err)
	require.NoError(t, s.SaveResource(ctx, "cluster-a", testResource("uid-1", "v1", "Pod", "default", "api", nil, nil)))
	require.NoError(t, s.Close())

	reopened, err := NewStorage(&Config{DataSourceName: dsn})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reopened.Close()) })

	res, err := reopened.Search(ctx, "SELECT * from resources WHERE name = 'api'", storage.SQLPatternType, nil)
	require.NoError(t, err)
	require.Equal(t, 1, res.Total)
	require.Equal(t, "api", res.Resources[0].Name)
}

func TestStorage_HealthAndCanceledRefresh(t *testing.T) {
	s := newTestStorage(t)
	require.NoError(t, s.CheckStorageHealth(context.Background()))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, s.Refresh(ctx), context.Canceled)
}

func TestStorage_ConcurrentReadWriteSmoke(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	errCh := make(chan error, 20)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			obj := testResource(fmt.Sprintf("uid-%d", i), "v1", "Pod", "default", fmt.Sprintf("pod-%02d", i), map[string]string{"app": "web"}, nil)
			if err := s.SaveResource(ctx, "cluster-a", obj); err != nil {
				errCh <- err
				return
			}
			_, err := s.SearchByTerms(ctx, map[string]any{"cluster": "cluster-a", "kind": "Pod"}, &storage.Pagination{Page: 1, PageSize: 5})
			if err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err)
	}
	count, err := s.CountResources(ctx)
	require.NoError(t, err)
	require.Equal(t, 10, count)
}

func TestStorage_MissingRowsUseStorageErrors(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	err := s.DeleteResource(ctx, "cluster-a", testDefaultPodLookup("missing"))
	require.True(t, errors.Is(err, storage.ErrResourceNotFound))

	err = s.DeleteResourceGroupRule(ctx, "missing")
	require.ErrorIs(t, err, storage.ErrResourceGroupRuleNotFound)

	_, err = s.ListResourceGroupsBy(ctx, "missing")
	require.ErrorIs(t, err, storage.ErrResourceGroupRuleNotFound)
}

func newTestStorage(t *testing.T) *Storage {
	t.Helper()
	s, err := NewStorage(&Config{DataSourceName: testDSN(t)})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, s.Close()) })
	return s
}

func testDSN(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "karpor-search.db")
}

func seedSearchResources(t *testing.T, s *Storage) {
	t.Helper()
	ctx := context.Background()
	resources := []struct {
		cluster     string
		uid         string
		apiVersion  string
		kind        string
		namespace   string
		name        string
		labels      map[string]string
		annotations map[string]string
	}{
		{
			cluster:    "cluster-a",
			uid:        "uid-api",
			apiVersion: "v1",
			kind:       "Pod",
			namespace:  "default",
			name:       "api",
			labels: map[string]string{
				"app":                    "api",
				"app.kubernetes.io/name": "demo",
				"env":                    "prod",
			},
			annotations: map[string]string{"owner": "team-a"},
		},
		{
			cluster:     "cluster-a",
			uid:         "uid-web",
			apiVersion:  "v1",
			kind:        "Pod",
			namespace:   "ops",
			name:        "web",
			labels:      map[string]string{"app": "web", "env": "dev"},
			annotations: map[string]string{"owner": "team-b"},
		},
		{
			cluster:     "cluster-a",
			uid:         "uid-svc",
			apiVersion:  "v1",
			kind:        "Service",
			namespace:   "default",
			name:        "web",
			labels:      map[string]string{"app": "web", "env": "prod"},
			annotations: map[string]string{"owner": "team-a"},
		},
		{
			cluster:     "cluster-b",
			uid:         "uid-db",
			apiVersion:  "apps/v1",
			kind:        "StatefulSet",
			namespace:   "data",
			name:        "db",
			labels:      map[string]string{"app": "db", "env": "prod"},
			annotations: map[string]string{"owner": "team-c"},
		},
	}
	for _, resource := range resources {
		require.NoError(t, s.SaveResource(ctx, resource.cluster, testResource(
			resource.uid,
			resource.apiVersion,
			resource.kind,
			resource.namespace,
			resource.name,
			resource.labels,
			resource.annotations,
		)))
	}
}

func testResource(uid, apiVersion, kind, namespace, name string, labels, annotations map[string]string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{Object: map[string]interface{}{
		"spec": map[string]interface{}{
			"name": name,
		},
	}}
	obj.SetAPIVersion(apiVersion)
	obj.SetKind(kind)
	obj.SetNamespace(namespace)
	obj.SetName(name)
	obj.SetUID(types.UID(uid))
	obj.SetLabels(labels)
	obj.SetAnnotations(annotations)
	obj.SetCreationTimestamp(metav1.Now())
	return obj
}

func testDefaultPodLookup(name string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("v1")
	obj.SetKind("Pod")
	obj.SetNamespace("default")
	obj.SetName(name)
	return obj
}

func resourceNames(resources []*storage.Resource) []string {
	names := make([]string, 0, len(resources))
	for _, resource := range resources {
		names = append(names, resource.Name)
	}
	return names
}
