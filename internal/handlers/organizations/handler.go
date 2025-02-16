package organizations

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/handlers/common"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-tfe"
)

func EntitlementsHandler(backendURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := common.GetTFEClientFromContext(r, backendURL)
		if err != nil {
			common.RespondWithError(w, err)
			return
		}

		vars := mux.Vars(r)
		orgName := vars["organization_name"]

		ctx, cancel := common.WithTimeout(r.Context(), common.DefaultTimeout)
		defer cancel()

		org, err := client.Organizations.Read(ctx, orgName)
		if err != nil {
			if err == tfe.ErrResourceNotFound {
				common.RespondWithError(w, common.HandlerError{
					Status:  http.StatusNotFound,
					Message: "Organization not found",
				})
				return
			}
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to fetch organization",
			})
			return
		}

		entitlements := &tfe.Entitlements{
			ID:                    org.Name,
			Operations:            true,
			PrivateModuleRegistry: true,
			Sentinel:              true,
			StateStorage:          true,
			Teams:                 true,
			VCSIntegrations:       true,
		}

		common.RespondWithJSON(w, http.StatusOK, entitlements)
	}
}
