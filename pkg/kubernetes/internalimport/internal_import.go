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

package internalimport

// To use k8s.io/kubernetes as a library, it is required to import the relevant packages related to k8s and specify the specific version in go.mod.
import (
	_ "k8s.io/api"
	_ "k8s.io/apimachinery"
	_ "k8s.io/apiserver"
	_ "k8s.io/cli-runtime"
	_ "k8s.io/client-go"
	_ "k8s.io/cloud-provider"
	_ "k8s.io/cluster-bootstrap"
	_ "k8s.io/code-generator"
	_ "k8s.io/component-base"
	_ "k8s.io/component-helpers"
	_ "k8s.io/controller-manager"
	_ "k8s.io/cri-api"
	_ "k8s.io/csi-translation-lib"
	_ "k8s.io/dynamic-resource-allocation"
	_ "k8s.io/kms"
	_ "k8s.io/kube-aggregator"
	_ "k8s.io/kube-controller-manager"
	_ "k8s.io/kube-proxy"
	_ "k8s.io/kube-scheduler"
	_ "k8s.io/kubectl"
	_ "k8s.io/kubelet"
	_ "k8s.io/legacy-cloud-providers"
	_ "k8s.io/metrics"
	_ "k8s.io/mount-utils"
	_ "k8s.io/pod-security-admission"
	_ "k8s.io/sample-apiserver"
)
