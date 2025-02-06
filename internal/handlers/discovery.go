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
			"login.v1":   baseURL + "/api/v1/",
			"modules.v1": baseURL + "/api/v1/",
			"state.v1":   baseURL + "/api/v1/",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(discoveryDoc)
	}
}
