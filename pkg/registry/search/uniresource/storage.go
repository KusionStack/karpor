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

package uniresource

import (
	"context"
	"fmt"

	"github.com/KusionStack/karbour/pkg/apis/search"
	"github.com/KusionStack/karbour/pkg/search/storage"
	filtersutil "github.com/KusionStack/karbour/pkg/util/filters"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	_ rest.Storage            = &REST{}
	_ rest.Scoper             = &REST{}
	_ rest.Lister             = &REST{}
	_ rest.ShortNamesProvider = &REST{}
)

type REST struct {
	Storage storage.SearchStorage
}

func NewREST(searchStorageGetter storage.SearchStorageGetter) (rest.Storage, error) {
	searchStorage, err := searchStorageGetter.GetSearchStorage()
	if err != nil {
		return nil, err
	}

	return &REST{
		Storage: searchStorage,
	}, nil
}

func (r *REST) New() runtime.Object {
	return &search.UniResource{}
}

func (r *REST) Destroy() {
}

func (r *REST) NamespaceScoped() bool {
	return false
}

func (r *REST) NewList() runtime.Object {
	return &search.UniResourceList{}
}

func (r *REST) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	queryString, ok := filtersutil.SearchQueryFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("query can't be empty")
	}

	patternType, ok := filtersutil.PatternTypeFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("pattern type can't be empty")
	}

	res, err := r.Storage.Search(ctx, queryString, patternType)
	if err != nil {
		return nil, err
	}

	rt := &search.UniResourceList{}
	for _, resource := range res.Resources {
		unObj := &unstructured.Unstructured{}
		unObj.SetUnstructuredContent(resource.Object)
		rt.Items = append(rt.Items, unObj)
	}
	return rt, nil
}

func (r *REST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: add real logic of convert to table when the storage layer is implemented
	return rest.NewDefaultTableConvertor(search.Resource("uniresources")).ConvertToTable(ctx, object, tableOptions)
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"ur"}
}
