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

package search

import (
	"encoding/json"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	podData = []byte(`
{
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
                "name": "pause1"
            },
            {
                "image": "registry.k8s.io/pause:3.8",
                "imagePullPolicy": "IfNotPresent",
                "name": "pause2"
            }
        ]
    }
}`)

	eventData = []byte(`
{
    "apiVersion": "v1",
	"count": 1,
    "involvedObject": {
        "apiVersion": "v1",
        "kind": "Pod",
        "name": "karpor-server-db4c78b4b-5jhhn",
        "namespace": "karpor"
    },
    "kind": "Event",
    "metadata": {
        "creationTimestamp": "2024-09-18T08:52:22Z",
        "name": "karpor-server-db4c78b4b-5jhhn.17f64a9c530cdb7c",
        "namespace": "karpor"
    },
    "reason": "Scheduled"
}`)

	pod      corev1.Pod
	event    corev1.Event
	unObjPod unstructured.Unstructured
)

func init() {
	json.Unmarshal(podData, &pod)
	json.Unmarshal(eventData, &event)
	json.Unmarshal(podData, &unObjPod)
}

func Test_customColumnsFormatter_Format(t *testing.T) {
	type input struct {
		spec string
		objs []runtime.Object
	}

	tests := []struct {
		name    string
		input   input
		want    any
		wantErr bool
	}{
		{name: "containers name", input: input{spec: "CONTAINER_NAME:spec.containers[*].name", objs: []runtime.Object{&unObjPod}}, want: CustomColumnsOutput{Fields: map[string]any{"CONTAINER_NAME": []any{"pause1", "pause2"}}, Titles: []string{"CONTAINER_NAME"}}, wantErr: false},
		{name: "dual containers name", input: input{spec: "CONTAINER_NAME:spec.containers[*].name", objs: []runtime.Object{&pod, &unObjPod}}, want: CustomColumnsOutput{Fields: map[string]any{"CONTAINER_NAME": []any{"pause1", "pause2"}}, Titles: []string{"CONTAINER_NAME"}}, wantErr: false},
		{name: "first container name", input: input{spec: "CONTAINER_NAME:spec.containers[0].name", objs: []runtime.Object{&unObjPod}}, want: CustomColumnsOutput{Fields: map[string]any{"CONTAINER_NAME": "pause1"}, Titles: []string{"CONTAINER_NAME"}}, wantErr: false},
		{name: "invalid path", input: input{spec: "NAME:{.kind} {.apiVersion}", objs: []runtime.Object{&unObjPod}}, want: nil, wantErr: true},
		{name: "name and apiVersion", input: input{spec: "NAME:metadata.name,API_VERSION:apiVersion", objs: []runtime.Object{&pod, &pod}}, want: CustomColumnsOutput{Fields: map[string]any{"NAME": "pause", "API_VERSION": "v1"}, Titles: []string{"NAME", "API_VERSION"}}, wantErr: false},
		{name: "count", input: input{spec: "COUNT:count", objs: []runtime.Object{&event}}, want: CustomColumnsOutput{Fields: map[string]any{"COUNT": int32(1)}, Titles: []string{"COUNT"}}, wantErr: false},
		{name: "not exist key", input: input{spec: "NOT_EXIST:spec.containers[*].xx.yy", objs: []runtime.Object{&event}}, want: CustomColumnsOutput{Fields: map[string]any{"NOT_EXIST": nil}, Titles: []string{"NOT_EXIST"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := NewCustomColumnsFormatter(tt.input.spec)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewCustomColumnsFormatter error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			for _, obj := range tt.input.objs {
				got, err := formatter.Format(obj)

				if (err != nil) != tt.wantErr {
					t.Errorf("customColumnsFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("customColumnsFormatter.Format() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestParseObjectFormatter(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		want    Formatter
		wantErr bool
	}{
		{name: "empty format", format: "", want: &NopFormatter{}, wantErr: false},
		{name: "origin format", format: "origin", want: &NopFormatter{}, wantErr: false},
		{name: "unsupported format", format: "yaml", want: nil, wantErr: true},
		{name: "empty custom-columns spec", format: "custom-columns=", want: &customColumnsFormatter{}, wantErr: true},
		{name: "custom-columns", format: "custom-columns=NAME:.metadata.name", want: &customColumnsFormatter{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseObjectFormatter(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseObjectFormatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("ParseObjectFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNopFormatter_Format(t *testing.T) {
	tests := []struct {
		name    string
		obj     runtime.Object
		want    any
		wantErr bool
	}{
		{name: "pod", obj: &pod, want: &pod, wantErr: false},
		{name: "event", obj: &event, want: &event, wantErr: false},
		{name: "nil", obj: nil, want: nil, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &NopFormatter{}
			got, err := f.Format(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("NopFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NopFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
