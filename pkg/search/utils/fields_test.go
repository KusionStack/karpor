// Copyright 2017 The Karbour Authors.
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

package utils

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestJSONPathFields(t *testing.T) {
	u := unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": metav1.ObjectMeta{
				Name: "foo",
				Labels: map[string]string{
					"label1":        "bar",
					"label2/label3": "foo",
				},
			},
		},
	}

	fields := NewJSONPathFields(NewJSONPathParser(), u.Object)

	testData := map[string]string{
		`{.metadata.name}`:                    "foo",
		`metadata.name`:                       "foo",
		`.metadata.name`:                      "foo",
		`{.metadata.labels.label1}`:           "bar",
		`{.metadata.labels['label1']}`:        "bar",
		`{.metadata.labels['label2/label3']}`: "foo",
		`{.notExistField}`:                    "",
	}

	for path, expectVal := range testData {
		actualVal := fields.Get(path)
		if actualVal != expectVal {
			t.Errorf(`the value of path '%s' is expected to be '%s', but got '%s'`, path, expectVal, actualVal)
		}
	}
}
