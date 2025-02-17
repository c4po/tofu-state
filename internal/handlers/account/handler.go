package account

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/handlers/common"
	"go.uber.org/zap"
)

func DetailsHandler(backendURL string, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Handling account details request")

		client, err := common.GetTFEClientFromContext(r, backendURL)
		if err != nil {
			logger.Error("Failed to get TFE client", zap.Error(err))
			common.RespondWithError(w, err)
			return
		}

		ctx, cancel := common.WithTimeout(r.Context(), common.DefaultTimeout)
		defer cancel()

		logger.Debug("Fetching current user details")
		user, err := client.Users.ReadCurrent(ctx)
		if err != nil {
			logger.Error("Failed to fetch user details", zap.Error(err))
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to fetch user details",
			})
			return
		}

		logger.Debug("Successfully fetched user details", zap.String("username", user.Username))
		common.RespondWithJSON(w, http.StatusOK, user)
	}
}
