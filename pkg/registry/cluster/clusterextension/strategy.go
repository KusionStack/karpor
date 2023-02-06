/*
Copyright The Alipay Authors.

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

package clusterextension

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

	"code.alipay.com/multi-cluster/karbour/pkg/apis/cluster"
)

// NewStrategy creates and returns a ClusterStrategy instance
func NewStrategy(typer runtime.ObjectTyper) ClusterStrategy {
	return ClusterStrategy{typer, names.SimpleNameGenerator}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a Fischer
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*cluster.ClusterExtension)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Fischer")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// MatchFischer is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchFischer(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *cluster.ClusterExtension) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type ClusterStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (ClusterStrategy) NamespaceScoped() bool {
	return false
}

func (ClusterStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (ClusterStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (ClusterStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (ClusterStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string { return nil }

func (ClusterStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (ClusterStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (ClusterStrategy) Canonicalize(obj runtime.Object) {
}

func (ClusterStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (ClusterStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
