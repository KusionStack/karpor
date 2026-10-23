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
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/KusionStack/karpor/pkg/core/entity"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/KusionStack/karpor/pkg/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	runtimejson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	_ "modernc.org/sqlite"
)

var (
	_ storage.Storage                  = &Storage{}
	_ storage.ResourceStorage          = &Storage{}
	_ storage.ResourceGroupRuleStorage = &Storage{}
	_ storage.SearchStorage            = &Storage{}
)

type Storage struct {
	db            *sql.DB
	objectEncoder runtime.Encoder
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func NewStorage(cfg *Config) (*Storage, error) {
	dataSourceName := defaultDataSourceName
	if cfg != nil && cfg.DataSourceName != "" {
		dataSourceName = cfg.DataSourceName
	}

	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	s := &Storage{
		db: db,
		objectEncoder: runtimejson.NewSerializerWithOptions(
			runtimejson.DefaultMetaFactory,
			scheme.Scheme,
			scheme.Scheme,
			runtimejson.SerializerOptions{Yaml: false, Pretty: true, Strict: true},
		),
	}
	if err := s.init(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) init(ctx context.Context) error {
	for _, stmt := range []string{
		`PRAGMA busy_timeout = 5000`,
		`PRAGMA foreign_keys = ON`,
		`PRAGMA journal_mode = WAL`,
	} {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	for _, stmt := range schemaStatements {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return s.createResourceGroupRuleIfNotExists(ctx, defaultResourceGroupRuleName)
}

func (s *Storage) createResourceGroupRuleIfNotExists(ctx context.Context, ruleName string) error {
	if _, err := s.GetResourceGroupRule(ctx, ruleName); err == nil {
		return nil
	} else if !errors.Is(err, storage.ErrResourceGroupRuleNotFound) {
		return err
	}

	now := metav1.Now()
	return s.SaveResourceGroupRule(ctx, &entity.ResourceGroupRule{
		ID:          entity.UUID(),
		Name:        ruleName,
		Description: fmt.Sprintf("Default resource group rule for %s", ruleName),
		Fields:      []string{ruleName},
		CreatedAt:   &now,
		UpdatedAt:   &now,
	})
}

func (s *Storage) SaveResource(ctx context.Context, cluster string, obj runtime.Object) error {
	doc, err := s.resourceDocument(cluster, obj)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `INSERT INTO resources (
		uid,
		cluster,
		api_version,
		kind,
		namespace,
		name,
		labels,
		annotations,
		creation_timestamp,
		deletion_timestamp,
		owner_references,
		resource_version,
		content,
		sync_at,
		deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	ON CONFLICT(uid) DO UPDATE SET
		cluster = excluded.cluster,
		api_version = excluded.api_version,
		kind = excluded.kind,
		namespace = excluded.namespace,
		name = excluded.name,
		labels = excluded.labels,
		annotations = excluded.annotations,
		creation_timestamp = excluded.creation_timestamp,
		deletion_timestamp = excluded.deletion_timestamp,
		owner_references = excluded.owner_references,
		resource_version = excluded.resource_version,
		content = excluded.content,
		sync_at = excluded.sync_at,
		deleted = 0`,
		doc.uid,
		doc.cluster,
		doc.apiVersion,
		doc.kind,
		doc.namespace,
		doc.name,
		doc.labels,
		doc.annotations,
		doc.creationTimestamp,
		doc.deletionTimestamp,
		doc.ownerReferences,
		doc.resourceVersion,
		doc.content,
		doc.syncAt,
	)
	return err
}

func (s *Storage) Refresh(ctx context.Context) error {
	return ctx.Err()
}

func (s *Storage) SoftDeleteResource(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, uid, err := s.findResource(ctx, cluster, obj)
	if err != nil {
		return err
	}
	if uid == "" {
		uid = string(unObj.GetUID())
	}
	if uid == "" {
		return fmt.Errorf("resource uid cannot be empty: %w", storage.ErrResourceNotFound)
	}

	res, err := s.db.ExecContext(ctx, `UPDATE resources SET deleted = 1, sync_at = ? WHERE uid = ?`, nowString(), uid)
	if err != nil {
		return err
	}
	return requireRowsAffected(res, storage.ErrResourceNotFound)
}

func (s *Storage) DeleteResource(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, uid, err := s.findResource(ctx, cluster, obj)
	if err != nil {
		return err
	}
	if uid == "" {
		uid = string(unObj.GetUID())
	}
	if uid == "" {
		return fmt.Errorf("resource uid cannot be empty: %w", storage.ErrResourceNotFound)
	}

	res, err := s.db.ExecContext(ctx, `DELETE FROM resources WHERE uid = ?`, uid)
	if err != nil {
		return err
	}
	return requireRowsAffected(res, storage.ErrResourceNotFound)
}

func (s *Storage) GetResource(ctx context.Context, cluster string, obj runtime.Object) error {
	unObj, _, err := s.findResource(ctx, cluster, obj)
	if err != nil {
		return err
	}
	obj.(*unstructured.Unstructured).Object = unObj.Object
	return nil
}

func (s *Storage) CountResources(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM resources`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Storage) DeleteAllResources(ctx context.Context, cluster string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM resources WHERE cluster = ?`, cluster)
	return err
}

func (s *Storage) CheckStorageHealth(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Storage) GetResourceGroupRule(ctx context.Context, name string) (*entity.ResourceGroupRule, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, name, description, fields, created_at, updated_at, deleted_at
		FROM resource_group_rules
		WHERE name = ?
		LIMIT 1`, name)

	rule, err := scanResourceGroupRule(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrResourceGroupRuleNotFound
		}
		return nil, err
	}
	return rule, nil
}

func (s *Storage) SaveResourceGroupRule(ctx context.Context, data *entity.ResourceGroupRule) error {
	if data == nil {
		return fmt.Errorf("resource group rule cannot be nil")
	}
	id := data.ID
	if id == "" {
		id = entity.UUID()
	}
	fields, err := json.Marshal(data.Fields)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `INSERT INTO resource_group_rules (
		id,
		name,
		description,
		fields,
		created_at,
		updated_at,
		deleted_at
	) VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		description = excluded.description,
		fields = excluded.fields,
		created_at = excluded.created_at,
		updated_at = excluded.updated_at,
		deleted_at = excluded.deleted_at`,
		id,
		data.Name,
		data.Description,
		string(fields),
		timePtrToString(data.CreatedAt),
		timePtrToString(data.UpdatedAt),
		timePtrToString(data.DeletedAt),
	)
	return err
}

func (s *Storage) DeleteResourceGroupRule(ctx context.Context, name string) error {
	rule, err := s.GetResourceGroupRule(ctx, name)
	if err != nil {
		return err
	}
	res, err := s.db.ExecContext(ctx, `DELETE FROM resource_group_rules WHERE id = ?`, rule.ID)
	if err != nil {
		return err
	}
	return requireRowsAffected(res, storage.ErrResourceGroupRuleNotFound)
}

func (s *Storage) ListResourceGroupRules(ctx context.Context) ([]*entity.ResourceGroupRule, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, description, fields, created_at, updated_at, deleted_at
		FROM resource_group_rules
		ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*entity.ResourceGroupRule
	for rows.Next() {
		rule, err := scanResourceGroupRule(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return nil, storage.ErrResourceGroupRuleNotFound
	}
	return rules, nil
}

func (s *Storage) CountResourceGroupRules(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM resource_group_rules`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Storage) ListResourceGroupsBy(ctx context.Context, ruleName string) (*storage.ResourceGroupResult, error) {
	rule, err := s.GetResourceGroupRule(ctx, ruleName)
	if err != nil {
		return nil, err
	}

	aggs, err := s.AggregateByTerms(ctx, rule.Fields)
	if err != nil {
		return nil, err
	}
	if aggs.Total == 0 {
		return nil, storage.ErrResourceGroupNotFound
	}

	groups := make([]*entity.ResourceGroup, 0, len(aggs.Buckets))
	for _, bucket := range aggs.Buckets {
		if len(rule.Fields) != len(bucket.Keys) {
			return nil, fmt.Errorf("mismatched number of fields: expected %d, got %d", len(rule.Fields), len(bucket.Keys))
		}
		rg := &entity.ResourceGroup{}
		for i, value := range bucket.Keys {
			assignResourceGroupField(rg, rule.Fields[i], value)
		}
		groups = append(groups, rg)
	}

	return &storage.ResourceGroupResult{
		Groups: groups,
		Fields: rule.Fields,
	}, nil
}

type resourceDocument struct {
	uid               string
	cluster           string
	apiVersion        string
	kind              string
	namespace         string
	name              string
	labels            string
	annotations       string
	creationTimestamp interface{}
	deletionTimestamp interface{}
	ownerReferences   string
	resourceVersion   string
	content           string
	syncAt            string
}

func (s *Storage) resourceDocument(cluster string, obj runtime.Object) (*resourceDocument, error) {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	uid := string(metaObj.GetUID())
	if uid == "" {
		return nil, fmt.Errorf("resource uid cannot be empty")
	}

	apiVersion, kind := resourceAPIVersionKind(obj)
	if apiVersion == "" || kind == "" {
		return nil, fmt.Errorf("resource apiVersion and kind cannot be empty")
	}

	content := bytes.NewBuffer(nil)
	if err := s.objectEncoder.Encode(obj, content); err != nil {
		return nil, err
	}
	labels, err := jsonString(metaObj.GetLabels(), "{}")
	if err != nil {
		return nil, err
	}
	annotations, err := jsonString(metaObj.GetAnnotations(), "{}")
	if err != nil {
		return nil, err
	}
	ownerReferences, err := jsonString(metaObj.GetOwnerReferences(), "[]")
	if err != nil {
		return nil, err
	}

	return &resourceDocument{
		uid:               uid,
		cluster:           cluster,
		apiVersion:        apiVersion,
		kind:              kind,
		namespace:         metaObj.GetNamespace(),
		name:              metaObj.GetName(),
		labels:            labels,
		annotations:       annotations,
		creationTimestamp: timeToString(metaObj.GetCreationTimestamp()),
		deletionTimestamp: timePtrToString(metaObj.GetDeletionTimestamp()),
		ownerReferences:   ownerReferences,
		resourceVersion:   metaObj.GetResourceVersion(),
		content:           content.String(),
		syncAt:            nowString(),
	}, nil
}

func (s *Storage) findResource(ctx context.Context, cluster string, obj runtime.Object) (*unstructured.Unstructured, string, error) {
	unObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, "", fmt.Errorf("only support *unstructured.Unstructured type")
	}
	apiVersion, kind := resourceAPIVersionKind(unObj)
	row := s.db.QueryRowContext(ctx, `SELECT uid, cluster, api_version, kind, namespace, name, content, sync_at, deleted
		FROM resources
		WHERE cluster = ? AND api_version = ? AND kind = ? AND namespace = ? AND name = ?
		ORDER BY deleted ASC, sync_at DESC, uid DESC
		LIMIT 1`,
		cluster,
		apiVersion,
		kind,
		unObj.GetNamespace(),
		unObj.GetName(),
	)
	res, uid, err := scanResource(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", fmt.Errorf("no resource found for cluster: %s, namespace: %s, name: %s: %w",
				cluster,
				unObj.GetNamespace(),
				unObj.GetName(),
				storage.ErrResourceNotFound,
			)
		}
		return nil, "", err
	}
	out := &unstructured.Unstructured{}
	out.SetUnstructuredContent(res.Object)
	return out, uid, nil
}

func resourceAPIVersionKind(obj runtime.Object) (string, string) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	apiVersion := gvk.GroupVersion().String()
	kind := gvk.Kind
	if unObj, ok := obj.(*unstructured.Unstructured); ok {
		if apiVersion == "" {
			apiVersion = unObj.GetAPIVersion()
		}
		if kind == "" {
			kind = unObj.GetKind()
		}
	}
	return apiVersion, kind
}

func scanResource(scanner rowScanner) (*storage.Resource, string, error) {
	var (
		uid        string
		cluster    string
		apiVersion string
		kind       string
		namespace  string
		name       string
		content    string
		syncAt     string
		deleted    int
	)
	if err := scanner.Scan(&uid, &cluster, &apiVersion, &kind, &namespace, &name, &content, &syncAt, &deleted); err != nil {
		return nil, "", err
	}
	res, err := storage.Map2Resource(map[string]interface{}{
		"cluster":    cluster,
		"apiVersion": apiVersion,
		"kind":       kind,
		"namespace":  namespace,
		"name":       name,
		"content":    content,
		"syncAt":     syncAt,
		"deleted":    deleted != 0,
	})
	if err != nil {
		return nil, "", err
	}
	return res, uid, nil
}

func scanResourceGroupRule(scanner rowScanner) (*entity.ResourceGroupRule, error) {
	var (
		id          string
		name        string
		description string
		fieldsJSON  string
		createdAt   sql.NullString
		updatedAt   sql.NullString
		deletedAt   sql.NullString
	)
	if err := scanner.Scan(&id, &name, &description, &fieldsJSON, &createdAt, &updatedAt, &deletedAt); err != nil {
		return nil, err
	}

	var fields []string
	if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
		return nil, err
	}

	return &entity.ResourceGroupRule{
		ID:          id,
		Name:        name,
		Description: description,
		Fields:      fields,
		CreatedAt:   stringToTime(createdAt),
		UpdatedAt:   stringToTime(updatedAt),
		DeletedAt:   stringToTime(deletedAt),
	}, nil
}

func requireRowsAffected(result sql.Result, notFound error) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return notFound
	}
	return nil
}

func timeToString(t metav1.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func timePtrToString(t *metav1.Time) interface{} {
	if t == nil || t.IsZero() {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}

func stringToTime(value sql.NullString) *metav1.Time {
	if !value.Valid || value.String == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, value.String)
	if err != nil {
		return nil
	}
	return &metav1.Time{Time: t}
}

func nowString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func jsonString(in interface{}, empty string) (string, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	if string(b) == "null" {
		return empty, nil
	}
	return string(b), nil
}

func assignResourceGroupField(rg *entity.ResourceGroup, field, value string) {
	switch field {
	case "cluster":
		rg.Cluster = value
	case "apiVersion":
		rg.APIVersion = value
	case "kind":
		rg.Kind = value
	case "namespace":
		rg.Namespace = value
	case "name":
		rg.Name = value
	default:
		switch {
		case strings.HasPrefix(field, "labels."):
			if rg.Labels == nil {
				rg.Labels = map[string]string{}
			}
			rg.Labels[strings.TrimPrefix(field, "labels.")] = value
		case strings.HasPrefix(field, "annotations."):
			if rg.Annotations == nil {
				rg.Annotations = map[string]string{}
			}
			rg.Annotations[strings.TrimPrefix(field, "annotations.")] = value
		}
	}
}
