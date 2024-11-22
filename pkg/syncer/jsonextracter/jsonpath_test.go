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

package jsonextracter

import (
	"encoding/json"
	"testing"
)

type jsonPathTest struct {
	name        string
	template    string
	input       map[string]interface{}
	expect      string
	expectError bool
}

func (t *jsonPathTest) Prepare(allowMissingKeys bool) (*JSONPath, error) {
	jp := New(t.name)
	jp.AllowMissingKeys(allowMissingKeys)
	return jp, jp.Parse(t.template)
}

func benchmarkJSONPath(test jsonPathTest, allowMissingKeys bool, b *testing.B) {
	jp, err := test.Prepare(allowMissingKeys)
	if err != nil {
		if !test.expectError {
			b.Errorf("in %s, parse %s error %v", test.name, test.template, err)
			return
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		jp.Extract(test.input)
	}
}

func testJSONPath(tests []jsonPathTest, allowMissingKeys bool, t *testing.T) {
	for _, test := range tests {
		jp, err := test.Prepare(allowMissingKeys)
		if err != nil {
			if !test.expectError {
				t.Errorf("in %s, parse %s error %v", test.name, test.template, err)
				continue
			}
		}

		got, err := jp.Extract(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf(`in %s, expected execute error, got %q`, test.name, got)
			}
		} else if err != nil {
			t.Errorf("in %s, execute error %v", test.name, err)
		}

		bytes_, _ := json.Marshal(got)
		out := string(bytes_)

		if out != test.expect {
			t.Errorf(`in %s, expect to get "%s", got "%s"`, test.name, test.expect, out)
		}
	}
}

var (
	pod = []byte(`{
		"apiVersion": "v1",
		"kind": "Pod",
		"metadata": {
			"labels": {
				"name": "pause",
				"app": "pause"
			},
			"name": "pause",
			"namespace": "default"
		},
		"spec": {
			"containers": [
				{
					"image": "registry.k8s.io/pause:3.8",
					"imagePullPolicy": "IfNotPresent",
					"name": "pause1",
					"resources": {
						"limits": {
							"cpu": "100m",
							"memory": "128Mi"
						},
						"requests": {
							"cpu": "100m",
							"memory": "128Mi"
						}
					}
				},
				{
					"image": "registry.k8s.io/pause:3.8",
					"imagePullPolicy": "IfNotPresent",
					"name": "pause2",
					"resources": {
						"limits": {
							"cpu": "10m",
							"memory": "64Mi"
						},
						"requests": {
							"cpu": "10m",
							"memory": "64Mi"
						}
					}
				}
			]
		}
}`)

	podData map[string]interface{}

	arbitrary = map[string]interface{}{
		"e": struct {
			F1 string `json:"f1"`
			F2 string `json:"f2"`
		}{F1: "f1", F2: "f2"},
	}
)

func init() {
	json.Unmarshal(pod, &podData)
}

func TestJSONPath(t *testing.T) {
	podTests := []jsonPathTest{
		{"empty", ``, podData, `null`, false},
		{"containers name", `{.kind}`, podData, `{"kind":"Pod"}`, false},
		{"containers name", `{.spec.containers[*].name}`, podData, `{"spec":{"containers":[{"name":"pause1"},{"name":"pause2"}]}}`, false},
		{"containers name (range)", `{range .spec.containers[*]}{.name}{end}`, podData, `{"spec":{"containers":[{"name":"pause1"},{"name":"pause2"}]}}`, false},
		{"containers name and image", `{.spec.containers[*]['name', 'image']}`, podData, `{"spec":{"containers":[{"image":"registry.k8s.io/pause:3.8","name":"pause1"},{"image":"registry.k8s.io/pause:3.8","name":"pause2"}]}}`, false},
		{"containers name and cpu", `{.spec.containers[*]['name', 'resources.requests.cpu']}`, podData, `{"spec":{"containers":[{"name":"pause1","resources":{"requests":{"cpu":"100m"}}},{"name":"pause2","resources":{"requests":{"cpu":"10m"}}}]}}`, false},
		{"container pause1 name and image", `{.spec.containers[?(@.name=="pause1")]['name', 'image']}`, podData, `{"spec":{"containers":[{"image":"registry.k8s.io/pause:3.8","name":"pause1"}]}}`, false},
		{"pick one label", `{.metadata.labels.name}`, podData, `{"metadata":{"labels":{"name":"pause"}}}`, false},
		{"not exist label", `{.metadata.labels.xx.dd}`, podData, `null`, true},
	}

	testJSONPath(podTests, false, t)

	allowMissingTests := []jsonPathTest{
		{"containers image", `{.spec.containers[*]['xname', 'image']}`, podData, `{"spec":{"containers":[{"image":"registry.k8s.io/pause:3.8"},{"image":"registry.k8s.io/pause:3.8"}]}}`, false},
		{"not exist key", `{.spec.containers[*]['name', 'xx.dd']}`, podData, `{"spec":{"containers":[{"name":"pause1"},{"name":"pause2"}]}}`, false},
		{"not exist label", `{.metadata.labels.xx.dd}`, podData, `{"metadata":{"labels":{}}}`, false},
	}

	testJSONPath(allowMissingTests, true, t)
}

func BenchmarkJSONPath(b *testing.B) {
	t := jsonPathTest{"range nodes capacity", `{.kind}`, podData, `{"kind":"Pod"}`, false}
	benchmarkJSONPath(t, true, b)
}
