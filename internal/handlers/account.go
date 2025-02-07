package handlers

import (
	"encoding/json"
	"net/http"

	tfe "github.com/hashicorp/go-tfe"
)

func AccountDetailsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user from context (set by JWT middleware)
		email := r.Context().Value("userEmail").(string)

		// Use tfe.User struct for the main resource object
		user := &tfe.User{
			ID:               "user-" + email,
			Username:         email,
			Email:            email,
			IsServiceAccount: false,
			AvatarURL:        "https://www.gravatar.com/avatar/00000000000000000000000000000000?d=mp",
			Permissions: &tfe.UserPermissions{
				CanCreateOrganizations: true,
				CanChangeEmail:         true,
				CanChangeUsername:      true,
			},
		}

		// Use standard JSON API response structure
		response := struct {
			Data *tfe.User `json:"data"`
		}{
			Data: user,
		}

		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(response)
	}
}
