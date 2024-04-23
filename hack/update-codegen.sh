#!/usr/bin/env bash

# Copyright The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
CODEGEN_PKG="${REPO_ROOT}/hack"
tmp_go_path="${REPO_ROOT}/_go"

cleanup() {
  rm -rf "${tmp_go_path}"
}
trap "cleanup" EXIT SIGINT
cleanup

source "${REPO_ROOT}"/hack/util.sh
util:create_gopath_tree "${REPO_ROOT}" "${tmp_go_path}"

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
bash "${CODEGEN_PKG}/generate-groups.sh" all \
  github.com/KusionStack/karbour/pkg/kubernetes/generated github.com/KusionStack/karbour/pkg/kubernetes/apis \
  "cluster:v1beta1 search:v1beta1" \
  --output-base "${tmp_go_path}" \
  --go-header-file "${REPO_ROOT}"/hack/boilerplate.go.txt

bash "${CODEGEN_PKG}/generate-internal-groups.sh" "deepcopy,defaulter,conversion,openapi" \
  github.com/KusionStack/karbour/pkg/kubernetes/generated github.com/KusionStack/karbour/pkg/kubernetes/apis github.com/KusionStack/karbour/pkg/kubernetes/apis \
  "cluster:v1beta1 search:v1beta1" \
  --output-base "${tmp_go_path}" \
  --go-header-file "${REPO_ROOT}/hack/boilerplate.go.txt"

# To use your own boilerplate text append:
#   --go-header-file "${REPO_ROOT}/hack/custom-boilerplate.go.txt"
