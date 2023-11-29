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

package config

import (
	"encoding/json"
	"net/http"

	"github.com/KusionStack/karbour/pkg/controller/config"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
)

func Get(configCtrl *config.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := ctxutil.GetLogger(r.Context())

		log.Info("Starting get config ...")

		b, err := json.MarshalIndent(configCtrl.Get(), "", "  ")
		if err != nil {
			log.Error(err, "Failed to mashal json")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(b)
	}
}
