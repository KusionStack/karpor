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
