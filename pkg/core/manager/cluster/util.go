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

package cluster

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// maskContent to apply MD5 hash and mask the content.
func maskContent(content string) string {
	// Apply MD5 hash
	hash := fmt.Sprintf("%x", md5.Sum([]byte(content)))

	// Calculate the range for masking
	maskLength := len(hash) * 3 / 4 // Three quarters of the hash length
	start := len(hash) / 8          // Start masking a quarter in
	end := start + maskLength       // End masking
	masked := hash[:start] + strings.Repeat("*", maskLength) + hash[end:]

	return masked
}
