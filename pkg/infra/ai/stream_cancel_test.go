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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// streamGoroutineRunning reports whether a goroutine started by a GenerateStream
// method is still alive.
func streamGoroutineRunning() bool {
	buf := make([]byte, 1<<20)
	buf = buf[:runtime.Stack(buf, true)]
	return strings.Contains(string(buf), ").GenerateStream.func")
}

// When the consumer stops reading and cancels the context, the goroutine that
// feeds the channel has to exit instead of blocking on a full channel.
func TestGenerateStream_CancelStopsGoroutine(t *testing.T) {
	clients := map[string]struct {
		client AIProvider
		cfg    AIConfig
	}{
		"openai":   {&OpenAIClient{}, AIConfig{AuthToken: "t", Model: "m"}},
		"azure":    {&AzureAIClient{}, AIConfig{AuthToken: "t", Model: "m"}},
		"deepseek": {&DeepseekClient{}, AIConfig{AuthToken: "t", Model: "m"}},
	}
	for name, tc := range clients {
		t.Run(name, func(t *testing.T) {
			release := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(http.StatusOK)
				// More chunks than the channel buffer holds, so a consumer that does
				// not read leaves the producer blocked on send.
				for i := 0; i < 300; i++ {
					fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":{\"content\":\"c%d\"}}]}\n\n", i)
				}
				w.(http.Flusher).Flush()
				select {
				case <-release:
				case <-r.Context().Done():
				}
			}))
			defer server.Close()
			defer close(release)

			tc.cfg.BaseURL = server.URL
			assert.NoError(t, tc.client.Configure(tc.cfg))

			ctx, cancel := context.WithCancel(context.Background())
			_, err := tc.client.GenerateStream(ctx, "prompt")
			assert.NoError(t, err)

			// Let the producer fill the channel, then stop consuming.
			time.Sleep(200 * time.Millisecond)
			cancel()

			deadline := time.Now().Add(3 * time.Second)
			for streamGoroutineRunning() && time.Now().Before(deadline) {
				time.Sleep(20 * time.Millisecond)
			}
			assert.False(t, streamGoroutineRunning(), "stream goroutine leaked after ctx was cancelled")
		})
	}
}
