package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func DiscoveryHandler(apiPath string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Handling discovery request")

		discoveryDoc := map[string]interface{}{
			"modules.v1": "/api/registry/v1/modules/",
			"state.v2":   apiPath,
			"tfe.v2":     apiPath,
			"tfe.v2.1":   apiPath,
			"tfe.v2.2":   apiPath,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(discoveryDoc); err != nil {
			logger.Error("Failed to encode discovery document", zap.Error(err))
		}
		logger.Debug("Discovery request completed successfully")
	}
}

// PingHandler handles the ping endpoint which is used for health checks
func PingHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Handling ping request")
		w.WriteHeader(http.StatusOK)
	}
}
