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
	"fmt"
	"strings"
	"testing"
)

func TestFieldPath(t *testing.T) {
	type args struct {
		obj              map[string]interface{}
		allowMissingKeys bool
		fields           []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{obj: podData, allowMissingKeys: true, fields: []string{}}, `null`, false},
		{"nil fields", args{obj: podData, allowMissingKeys: true, fields: nil}, `null`, false},
		{"nil input nil fields", args{obj: nil, allowMissingKeys: true, fields: nil}, `null`, false},
		{"nil input non-nil fields", args{obj: nil, allowMissingKeys: true, fields: []string{"xx"}}, `{}`, false},
		{"nil input non-nil fields not allow missing", args{obj: nil, allowMissingKeys: false, fields: []string{"xx"}}, `null`, true},

		{"kind", args{obj: podData, allowMissingKeys: true, fields: []string{"kind"}}, `{"kind":"Pod"}`, false},
		{"lables", args{obj: podData, allowMissingKeys: true, fields: []string{"metadata", "labels"}}, `{"metadata":{"labels":{"app":"pause","name":"pause"}}}`, false},
		{"label name", args{obj: podData, allowMissingKeys: true, fields: []string{"metadata", "labels", "name"}}, `{"metadata":{"labels":{"name":"pause"}}}`, false},
		{"containers", args{obj: podData, allowMissingKeys: true, fields: []string{"spec", "containers"}}, `{"spec":{"containers":[{"image":"registry.k8s.io/pause:3.8","imagePullPolicy":"IfNotPresent","name":"pause1","resources":{"limits":{"cpu":"100m","memory":"128Mi"},"requests":{"cpu":"100m","memory":"128Mi"}}},{"image":"registry.k8s.io/pause:3.8","imagePullPolicy":"IfNotPresent","name":"pause2","resources":{"limits":{"cpu":"10m","memory":"64Mi"},"requests":{"cpu":"10m","memory":"64Mi"}}}]}}`, false},
		{"wrong type", args{obj: podData, allowMissingKeys: true, fields: []string{"metadata", "labels", "name", "xx"}}, "null", true},
		{"not allow miss key", args{obj: podData, allowMissingKeys: false, fields: []string{"metadata", "labels", "xx"}}, "null", true},
		{"allow miss key", args{obj: podData, allowMissingKeys: true, fields: []string{"metadata", "labels", "xx"}}, `{"metadata":{"labels":{}}}`, false},
		{"arbitrary", args{obj: arbitrary, allowMissingKeys: true, fields: []string{"e"}}, `{"e":{"f1":"f1","f2":"f2"}}`, false},
		{"not map", args{obj: arbitrary, allowMissingKeys: true, fields: []string{"e", "f1"}}, `null`, true},
	}

	fieldPathToJSONPath := func(nestedField []string) string {
		if nestedField == nil {
			return ""
		}
		if len(nestedField) == 0 {
			return "{}"
		}

		return fmt.Sprintf("{.%s}", strings.Join(nestedField, "."))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NestedFieldNoCopy(tt.args.obj, tt.args.allowMissingKeys, tt.args.fields...)

			jpt := jsonPathTest{tt.name, fieldPathToJSONPath(tt.args.fields), tt.args.obj, tt.want, tt.wantErr}
			testJSONPath(t, []jsonPathTest{jpt}, tt.args.allowMissingKeys)

			if (err != nil) != tt.wantErr {
				t.Errorf("NestedFieldNoCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			data, _ := json.Marshal(got)
			if string(data) != tt.want {
				t.Errorf("NestedFieldNoCopy() = %v, want %v", string(data), tt.want)
			}
		})
	}
}

func BenchmarkFieldPath(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NestedFieldNoCopy(podData, false, "kind")
	}
}
