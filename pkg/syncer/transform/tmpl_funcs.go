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

package transform

import (
	"sync"
	"text/template"
)

var (
	clusterTmplFuncs     = make(map[string]template.FuncMap)
	clusterTmplFuncsLock sync.RWMutex
)

func RegisterClusterTmplFunc(cluster, FuncName string, tmplFunc any) error {
	clusterTmplFuncsLock.Lock()
	defer clusterTmplFuncsLock.Unlock()

	// TODO: check whether the tmplFunc is valid as a template function.

	if _, ok := clusterTmplFuncs[cluster]; !ok {
		clusterTmplFuncs[cluster] = make(template.FuncMap)
	}

	clusterTmplFuncs[cluster][FuncName] = tmplFunc
	return nil
}

func GetClusterTmplFuncs(cluster string) (template.FuncMap, bool) {
	clusterTmplFuncsLock.RLock()
	defer clusterTmplFuncsLock.RUnlock()

	funcs, found := clusterTmplFuncs[cluster]
	return funcs, found
}
