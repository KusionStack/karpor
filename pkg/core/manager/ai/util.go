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

package ai

import (
	"regexp"
	"strings"
)

// IsInvalidQuery check if the query is invalid
func IsInvalidQuery(sql string) bool {
	return strings.Contains(strings.ToLower(sql), "error")
}

// ExtractSelectSQL extracts SQL statements that start with "SELECT * FROM"
func ExtractSelectSQL(sql string) string {
	res := regexp.MustCompile(`(?i)SELECT \* FROM [^;]+`)
	match := res.FindString(sql)
	return match
}
