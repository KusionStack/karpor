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
	"reflect"
	"testing"
)

func TestBuildExtracter(t *testing.T) {
	type args struct {
		path             string
		allowMissingKeys bool
	}
	tests := []struct {
		name    string
		args    args
		want    Extracter
		wantErr bool
	}{
		{name: "invalid path", args: args{path: `{`, allowMissingKeys: false}, want: nil, wantErr: true},
		{name: "fieldPath extracter", args: args{path: `{}`, allowMissingKeys: false}, want: &NestedFieldPath{}, wantErr: false},
		{name: "fieldPath extracter", args: args{path: ``, allowMissingKeys: false}, want: &NestedFieldPath{}, wantErr: false},
		{name: "fieldPath extracter", args: args{path: `{.metadata.labels.name}`, allowMissingKeys: false}, want: &NestedFieldPath{}, wantErr: false},
		{name: "fieldPath extracter", args: args{path: `{.metadata.labels['name']}`, allowMissingKeys: false}, want: &NestedFieldPath{}, wantErr: false},
		{name: "jsonPath extracter", args: args{path: `{.metadata.labels.name}{.metadata.labels.app}`, allowMissingKeys: false}, want: &JSONPath{}, wantErr: false},
		{name: "jsonPath extracter", args: args{path: `{.metadata.labels['name', 'app']}`, allowMissingKeys: false}, want: &JSONPath{}, wantErr: false},
		{name: "jsonPath extracter", args: args{path: `{.spec.containers[*].name}`, allowMissingKeys: false}, want: &JSONPath{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildExtracter(tt.args.path, tt.args.allowMissingKeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildExtracter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if reflect.TypeOf(tt.want) != reflect.TypeOf(got) {
				t.Errorf("BuildExtracter() = %T, want %T", got, tt.want)
			}
		})
	}
}
