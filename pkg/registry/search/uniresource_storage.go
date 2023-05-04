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

package search

import (
	"context"

	"github.com/KusionStack/karbour/pkg/apis/search"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	_ rest.Storage            = &UniResourceREST{}
	_ rest.Scoper             = &UniResourceREST{}
	_ rest.Lister             = &UniResourceREST{}
	_ rest.ShortNamesProvider = &UniResourceREST{}
)

type UniResourceREST struct{}

func NewUniResourceREST() rest.Storage {
	return &UniResourceREST{}
}

func (r *UniResourceREST) New() runtime.Object {
	return &search.UniResource{}
}

func (r *UniResourceREST) Destroy() {
}

func (r *UniResourceREST) NamespaceScoped() bool {
	return false
}

func (r *UniResourceREST) NewList() runtime.Object {
	return &search.UniResourceList{}
}

func (r *UniResourceREST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	// TODO: add real logic of list when the storage layer is implemented
	return &search.UniResourceList{}, nil
}

func (r *UniResourceREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: add real logic of convert to table when the storage layer is implemented
	return rest.NewDefaultTableConvertor(search.Resource("uniresources")).ConvertToTable(ctx, object, tableOptions)
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *UniResourceREST) ShortNames() []string {
	return []string{"ur"}
}
