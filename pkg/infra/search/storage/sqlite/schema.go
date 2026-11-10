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

package sqlite

const (
	defaultDataSourceName = "file:karpor-search.db?mode=rwc"

	defaultResourceGroupRuleName = "namespace"
)

var schemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS resources (
		uid TEXT PRIMARY KEY,
		cluster TEXT NOT NULL,
		api_version TEXT NOT NULL,
		kind TEXT NOT NULL,
		namespace TEXT NOT NULL DEFAULT '',
		name TEXT NOT NULL,
		labels TEXT NOT NULL DEFAULT '{}',
		annotations TEXT NOT NULL DEFAULT '{}',
		creation_timestamp TEXT,
		deletion_timestamp TEXT,
		owner_references TEXT NOT NULL DEFAULT '[]',
		resource_version TEXT,
		content TEXT NOT NULL,
		sync_at TEXT NOT NULL,
		deleted INTEGER NOT NULL DEFAULT 0 CHECK (deleted IN (0, 1))
	)`,
	`CREATE INDEX IF NOT EXISTS idx_resources_identity ON resources (
		cluster,
		api_version,
		kind,
		namespace,
		name
	)`,
	`CREATE INDEX IF NOT EXISTS idx_resources_cluster ON resources (cluster)`,
	`CREATE INDEX IF NOT EXISTS idx_resources_common_search ON resources (
		cluster,
		kind,
		namespace,
		name,
		deleted
	)`,
	`CREATE INDEX IF NOT EXISTS idx_resources_deleted_sync_at ON resources (
		deleted,
		sync_at
	)`,
	`CREATE TABLE IF NOT EXISTS resource_group_rules (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL DEFAULT '',
		fields TEXT NOT NULL DEFAULT '[]',
		created_at TEXT,
		updated_at TEXT,
		deleted_at TEXT
	)`,
	`CREATE INDEX IF NOT EXISTS idx_resource_group_rules_name ON resource_group_rules (name)`,
}
