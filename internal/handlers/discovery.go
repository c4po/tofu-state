package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/c4po/tofu-state/internal/config"
)

func DiscoveryHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		discoveryDoc := map[string]interface{}{
			"modules.v1": "/api/registry/v1/modules/",
			"state.v2":   "/api/v2/",
			"tfe.v2":     "/api/v2/",
			"tfe.v2.1":   "/api/v2/",
			"tfe.v2.2":   "/api/v2/",
		}

		// Add login.v1 only if OIDC is configured
		if cfg.OIDCConfig.ClientID != "" {
			discoveryDoc["login.v1"] = map[string]interface{}{
				"client":      "tofu-cli",
				"grant_types": []string{"authz_code"},
				"authz":       "/login",
				"token":       "/api/v1/login/token",
				"ports":       []int{10000, 10010},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(discoveryDoc)
	}
}
