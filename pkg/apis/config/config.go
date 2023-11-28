package config

import (
	"encoding/json"
	"net/http"

	"github.com/KusionStack/karbour/pkg/controller/config"
	"github.com/KusionStack/karbour/pkg/middleware"
	"k8s.io/klog/v2"
)

func Get(configCtrl *config.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, ok := r.Context().Value(middleware.APILoggerKey).(klog.Logger)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.Info("Starting get config ...")

		b, err := json.MarshalIndent(configCtrl.Get(), "", "  ")
		if err != nil {
			logger.Error(err, "Failed to mashal json")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(b)
	}
}
