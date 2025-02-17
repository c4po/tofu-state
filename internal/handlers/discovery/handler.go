package discovery

import (
	"encoding/json"
	"net/http"
)

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		discoveryDoc := map[string]interface{}{
			"modules.v1": "/api/registry/v1/modules/",
			"state.v2":   "/api/v2/",
			"tfe.v2":     "/api/tfe/v2/",
			"tfe.v2.1":   "/api/tfe/v2/",
			"tfe.v2.2":   "/api/tfe/v2/",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(discoveryDoc)
	}
}
