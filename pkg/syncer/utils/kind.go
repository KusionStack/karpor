// Copyright The Karpor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

// ResourceKind returns the kind used to query the stored resources of a sync rule.
// kind is the Kind reported by API discovery for the resource; when it is unknown,
// the resource name with its last letter removed is used, which only works for
// plurals that are the kind plus an "s" (e.g. "pods").
func ResourceKind(resource, kind string) string {
	if kind != "" {
		return kind
	}
	if resource == "" {
		return ""
	}
	return resource[0 : len(resource)-1]
}
