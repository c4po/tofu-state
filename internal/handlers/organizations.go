package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-tfe"
)

func OrganizationEntitlementsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgName := vars["organization_name"]

		entitlements := &tfe.Entitlements{
			ID:             orgName,
			StateStorage:   true,
			AuditLogging:   true,
			Agents:         false,
			CostEstimation: false,
		}

		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": entitlements,
		})
	}
}
