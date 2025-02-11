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

package safeutil

import (
	"testing"
	"time"

	"github.com/elliotxx/safe"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2/klogr"
)

func TestGo(t *testing.T) {
	type args struct {
		do safe.DoFunc
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful-recover-crash-in-safe-Go",
			args: args{
				do: func() {
					panic("ah, I'm down")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Go(tt.args.do)
			time.Sleep(100 * time.Microsecond)
		})
	}
}

func TestGoL(t *testing.T) {
	type args struct {
		do     safe.DoFunc
		logger logr.Logger
	}

	getTestingLogger := func() logr.Logger {
		logger := klogr.New()

		return logger
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful-recover-crash-in-safe-Go",
			args: args{
				do: func() {
					panic("ah, I'm down")
				},
				logger: getTestingLogger(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GoL(tt.args.do, tt.args.logger)
			time.Sleep(time.Second * 1)
		})
	}
}
