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

package elasticsearch

import (
	"encoding/json"

	"github.com/aquasecurity/esquery"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func generateIndexRequest(cluster string, obj runtime.Object) (id string, body []byte, err error) {
	metaObj, err := meta.Accessor(obj)
	if err != nil {
		return
	}

	body, err = json.Marshal(map[string]interface{}{
		apiVersionKey: obj.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		kindKey:       obj.GetObjectKind().GroupVersionKind().Kind,
		nameKey:       metaObj.GetName(),
		namespaceKey:  metaObj.GetNamespace(),
		clusterKey:    cluster,
		objectKey:     metaObj,
	})
	if err != nil {
		return
	}
	id = string(metaObj.GetUID())
	return
}

func generateQuery(cluster, namespace, name string, obj runtime.Object) map[string]interface{} {
	query := make(map[string]interface{})
	query["query"] = esquery.Bool().Must(
		esquery.Term(apiVersionKey, obj.GetObjectKind().GroupVersionKind().GroupVersion().String()),
		esquery.Term(kindKey, obj.GetObjectKind().GroupVersionKind().Kind),
		esquery.Term(nameKey, name),
		esquery.Term(namespaceKey, namespace),
		esquery.Term(clusterKey, cluster),
	).Map()
	return query
}
