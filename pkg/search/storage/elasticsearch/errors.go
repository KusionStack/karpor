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

package elasticsearch

import "fmt"

// ESError is an error type which represents a single ES error
type ESError struct {
	StatusCode int
	Message    string
}

func (e *ESError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.StatusCode, e.Message)
}
