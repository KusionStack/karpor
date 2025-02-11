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

package certgenerator

import (
	"context"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type FakeCoreV1 struct {
	v1.CoreV1Interface
}

func (FakeCoreV1) Secrets(namespace string) v1.SecretInterface {
	return &FakeSecret{}
}

func (FakeCoreV1) ConfigMaps(namespace string) v1.ConfigMapInterface {
	return &FakeConfigMap{}
}

type FakeSecret struct {
	v1.SecretInterface
}

func (f *FakeSecret) Apply(ctx context.Context, secret *applycorev1.SecretApplyConfiguration, opts metav1.ApplyOptions) (result *coreV1.Secret, err error) {
	return &coreV1.Secret{}, nil
}

type FakeConfigMap struct {
	v1.ConfigMapInterface
}

func (f *FakeConfigMap) Apply(ctx context.Context, secret *applycorev1.ConfigMapApplyConfiguration, opts metav1.ApplyOptions) (result *coreV1.ConfigMap, err error) {
	return &coreV1.ConfigMap{}, nil
}
