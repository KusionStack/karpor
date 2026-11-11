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

package aggregator

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// TestBuildPodLogOptions tests the buildPodLogOptions function.
func TestBuildPodLogOptions(t *testing.T) {
	sinceTime := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name        string
		query       url.Values
		download    bool
		expected    *corev1.PodLogOptions
		expectError string
	}{
		{
			name:  "Defaults for streaming",
			query: url.Values{},
			expected: &corev1.PodLogOptions{
				Follow:    true,
				TailLines: pointer.Int64(1000),
			},
		},
		{
			name:     "Defaults for download",
			query:    url.Values{},
			download: true,
			expected: &corev1.PodLogOptions{},
		},
		{
			name: "Container and timestamps",
			query: url.Values{
				"container":  {"nginx"},
				"timestamps": {"true"},
			},
			expected: &corev1.PodLogOptions{
				Container:  "nginx",
				Follow:     true,
				Timestamps: true,
				TailLines:  pointer.Int64(1000),
			},
		},
		{
			name:  "Valid since",
			query: url.Values{"since": {"1h30m"}},
			expected: &corev1.PodLogOptions{
				Follow:       true,
				SinceSeconds: pointer.Int64(5400),
				TailLines:    pointer.Int64(1000),
			},
		},
		{
			name:  "Sub-second since is rounded up",
			query: url.Values{"since": {"500ms"}},
			expected: &corev1.PodLogOptions{
				Follow:       true,
				SinceSeconds: pointer.Int64(1),
				TailLines:    pointer.Int64(1000),
			},
		},
		{
			name:        "Since with unsupported day unit",
			query:       url.Values{"since": {"2d"}},
			expectError: `invalid since "2d"`,
		},
		{
			name:        "Zero since",
			query:       url.Values{"since": {"0s"}},
			expectError: `invalid since "0s"`,
		},
		{
			name:        "Negative since",
			query:       url.Values{"since": {"-5m"}},
			expectError: `invalid since "-5m"`,
		},
		{
			name:  "Valid sinceTime",
			query: url.Values{"sinceTime": {"2024-01-02T03:04:05Z"}},
			expected: &corev1.PodLogOptions{
				Follow:    true,
				SinceTime: &metav1.Time{Time: sinceTime},
				TailLines: pointer.Int64(1000),
			},
		},
		{
			name:  "SinceTime with fractional seconds",
			query: url.Values{"sinceTime": {"2024-01-02T03:04:05.123Z"}},
			expected: &corev1.PodLogOptions{
				Follow:    true,
				SinceTime: &metav1.Time{Time: sinceTime.Add(123 * time.Millisecond)},
				TailLines: pointer.Int64(1000),
			},
		},
		{
			name:        "Invalid sinceTime",
			query:       url.Values{"sinceTime": {"yesterday"}},
			expectError: `invalid sinceTime "yesterday"`,
		},
		{
			name: "Since takes precedence over sinceTime",
			query: url.Values{
				"since":     {"10m"},
				"sinceTime": {"2024-01-02T03:04:05Z"},
			},
			expected: &corev1.PodLogOptions{
				Follow:       true,
				SinceSeconds: pointer.Int64(600),
				TailLines:    pointer.Int64(1000),
			},
		},
		{
			name:  "Valid tailLines",
			query: url.Values{"tailLines": {"50"}},
			expected: &corev1.PodLogOptions{
				Follow:    true,
				TailLines: pointer.Int64(50),
			},
		},
		{
			name:     "TailLines for download",
			query:    url.Values{"tailLines": {"50"}},
			download: true,
			expected: &corev1.PodLogOptions{
				TailLines: pointer.Int64(50),
			},
		},
		{
			name:        "Non-numeric tailLines",
			query:       url.Values{"tailLines": {"all"}},
			expectError: `invalid tailLines "all"`,
		},
		{
			name:        "Negative tailLines",
			query:       url.Values{"tailLines": {"-1"}},
			expectError: `invalid tailLines "-1"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := buildPodLogOptions(tt.query, tt.download)
			if tt.expectError != "" {
				require.ErrorContains(t, err, tt.expectError)
				require.Nil(t, opts)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, opts)
		})
	}
}
