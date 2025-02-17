package discovery

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func Handler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Handling discovery request")

		discoveryDoc := map[string]interface{}{
			"modules.v1": "/api/registry/v1/modules/",
			"state.v2":   "/api/v2/",
			"tfe.v2":     "/api/tfe/v2/",
			"tfe.v2.1":   "/api/tfe/v2/",
			"tfe.v2.2":   "/api/tfe/v2/",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(discoveryDoc); err != nil {
			logger.Error("Failed to encode discovery document", zap.Error(err))
		}
		logger.Debug("Discovery request completed successfully")
	}
}
