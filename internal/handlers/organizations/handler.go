package organizations

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/handlers/common"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-tfe"
	"go.uber.org/zap"
)

func EntitlementsHandler(backendURL string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Handling organization entitlements request")

		client, err := common.GetTFEClientFromContext(r, backendURL)
		if err != nil {
			logger.Error("Failed to get TFE client", zap.Error(err))
			common.RespondWithError(w, err)
			return
		}

		vars := mux.Vars(r)
		orgName := vars["organization_name"]
		logger.Debug("Fetching organization details", zap.String("organization", orgName))

		ctx, cancel := common.WithTimeout(r.Context(), common.DefaultTimeout)
		defer cancel()

		org, err := client.Organizations.Read(ctx, orgName)
		if err != nil {
			if err == tfe.ErrResourceNotFound {
				logger.Debug("Organization not found", zap.String("organization", orgName))
				common.RespondWithError(w, common.HandlerError{
					Status:  http.StatusNotFound,
					Message: "Organization not found",
				})
				return
			}
			logger.Error("Failed to fetch organization", zap.String("organization", orgName), zap.Error(err))
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to fetch organization",
			})
			return
		}

		logger.Debug("Successfully fetched organization", zap.String("organization", org.Name))

		entitlements := &tfe.Entitlements{
			ID:                    org.Name,
			Operations:            true,
			PrivateModuleRegistry: true,
			Sentinel:              true,
			StateStorage:          true,
			Teams:                 true,
			VCSIntegrations:       true,
		}

		logger.Debug("Responding with entitlements", zap.String("organization", org.Name))
		common.RespondWithJSON(w, http.StatusOK, entitlements)
	}
}
