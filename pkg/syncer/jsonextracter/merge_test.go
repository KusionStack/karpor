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
	"bytes"
	"encoding/json"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func BenchmarkJSONPathMerge(b *testing.B) {
	tests := []jsonPathTest{
		{"kind", `{.kind}`, podData, "", false},
		{"apiVersion", "{.apiVersion}", podData, "", false},
		{"metadata", "{.metadata}", podData, "", false},
	}

	extracters := make([]Extracter, 0)

	for _, test := range tests {
		ex, err := test.Prepare(false)
		if err != nil {
			if !test.expectError {
				b.Errorf("in %s, parse %s error %v", test.name, test.template, err)
			}
			return
		}
		extracters = append(extracters, ex)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		Merge(extracters, podData)
	}
}

func BenchmarkFieldPathMerge(b *testing.B) {
	fields := []string{"kind", "apiVersion", "metadata"}

	extracters := make([]Extracter, 0)

	for _, f := range fields {
		extracters = append(extracters, NewNestedFieldPath([]string{f}, false))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		Merge(extracters, podData)
	}
}

func BenchmarkTmpl(b *testing.B) {
	tmpl := `{"kind": "{{ .Object.kind }}","apiVersion": "{{ .Object.apiVersion}}","metadata": {{ toJson .Object.metadata }}}`
	obj := unstructured.Unstructured{Object: podData}

	t, _ := template.New("transformTemplate").Funcs(sprig.FuncMap()).Parse(tmpl)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		t.Execute(&buf, obj)

		var dest unstructured.Unstructured
		json.Unmarshal(buf.Bytes(), &dest)
	}
}

func TestMerge(t *testing.T) {
	containerName := jsonPathTest{"containers name", `{.spec.containers[*].name}`, podData, "", false}
	containerNameExtracter, _ := containerName.Prepare(true)

	containerImage := jsonPathTest{"containers image", `{.spec.containers[*].image}`, podData, "", false}
	containerImageExtracter, _ := containerImage.Prepare(true)

	kindExtracter := NewNestedFieldPath([]string{"kind"}, true)

	apiVersionExtracter := NewNestedFieldPath([]string{"apiVersion"}, true)

	type args struct {
		extracters []Extracter
		input      map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "merge name and image", args: args{extracters: []Extracter{containerImageExtracter, containerNameExtracter}, input: podData},
			want: `{"spec":{"containers":[{"name":"pause1"},{"name":"pause2"}]}}`, wantErr: false,
		},
		{
			name: "name kind apiVersion", args: args{extracters: []Extracter{containerNameExtracter, kindExtracter, apiVersionExtracter}, input: podData},
			want: `{"apiVersion":"v1","kind":"Pod","spec":{"containers":[{"name":"pause1"},{"name":"pause2"}]}}`, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Merge(tt.args.extracters, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			data, _ := json.Marshal(got)
			if string(data) != tt.want {
				t.Errorf("Merge() = %v, want %v", string(data), tt.want)
			}
		})
	}
}
