package account

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/handlers/common"
)

func DetailsHandler(backendURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := common.GetTFEClientFromContext(r, backendURL)
		if err != nil {
			common.RespondWithError(w, err)
			return
		}

		ctx, cancel := common.WithTimeout(r.Context(), common.DefaultTimeout)
		defer cancel()

		user, err := client.Users.ReadCurrent(ctx)
		if err != nil {
			common.RespondWithError(w, common.HandlerError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to fetch user details",
			})
			return
		}

		common.RespondWithJSON(w, http.StatusOK, user)
	}
}
