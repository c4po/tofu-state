package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/jsonapi"
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
			// ... initialize other fields as needed ...
		}

		w.Header().Set("Content-Type", "application/vnd.api+json")
		jsonapi.MarshalPayload(w, entitlements)
	}
}
