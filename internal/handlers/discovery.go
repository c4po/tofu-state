package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/c4po/tofu-state/internal/config"
)

func DiscoveryHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := "https://" + cfg.ExternalHost
		if cfg.CertFile == "" {
			baseURL = "http://" + cfg.ExternalHost
		}

		discoveryDoc := map[string]string{
			"login.v1":   baseURL + "/api/v1/login",
			"modules.v1": baseURL + "/api/v1/modules",
			"state.v1":   baseURL + "/api/v1/state",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(discoveryDoc)
	}
}
