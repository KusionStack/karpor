/*
Copyright The Karbour Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	"github.com/KusionStack/karbour/pkg/apis/cluster"
	"github.com/KusionStack/karbour/pkg/scheme"
)

var Strategy = clusterStrategy{scheme.Scheme, names.SimpleNameGenerator}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a Fischer
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*cluster.Cluster)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Fischer")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// MatchCluster is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchCluster(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *cluster.Cluster) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, false)
}

type clusterStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (clusterStrategy) NamespaceScoped() bool {
	return false
}

func (clusterStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	obj.(*cluster.Cluster).Status = cluster.ClusterStatus{}
}

func (clusterStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	obj.(*cluster.Cluster).Status = old.(*cluster.Cluster).Status
}

func (clusterStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (clusterStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string { return nil }

func (clusterStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (clusterStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (clusterStrategy) Canonicalize(obj runtime.Object) {
}

func (clusterStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (clusterStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

type clusterStatusStrategy struct {
	clusterStrategy
}

var StatusStartegy = clusterStatusStrategy{Strategy}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update of status
func (clusterStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	obj.(*cluster.Cluster).Spec = old.(*cluster.Cluster).Spec
}

// ValidateUpdate is the default update validation for an end user updating status
func (clusterStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (clusterStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
