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

package config

import _ "embed"

var DefaultConfig = [][]byte{DefaultSyncStrategy, DefaultAnonymousRBAC, DefaultGuestRBAC, DefaultAdminRBAC}

//go:embed default-anonymous-rbac.yaml
var DefaultAnonymousRBAC []byte

//go:embed default-karpor-admin-rbac.yaml
var DefaultAdminRBAC []byte

//go:embed default-karpor-guest-rbac.yaml
var DefaultGuestRBAC []byte

//go:embed default-relationship.yaml
var DefaultRelationship []byte

//go:embed default-sync-strategy.yaml
var DefaultSyncStrategy []byte

var CrdList = [][]byte{ClustersCrd, SyncRegistriesCrd, SyncResourcesCrd, TransformRulesCrd, TrimRulesCrd}

//go:embed crds/cluster.karpor.io_clusters.yaml
var ClustersCrd []byte

//go:embed crds/search.karpor.io_syncregistries.yaml
var SyncRegistriesCrd []byte

//go:embed crds/search.karpor.io_syncresources.yaml
var SyncResourcesCrd []byte

//go:embed crds/search.karpor.io_transformrules.yaml
var TransformRulesCrd []byte

//go:embed crds/search.karpor.io_trimrules.yaml
var TrimRulesCrd []byte

//go:embed agent.tpl
var AgentTpl []byte
