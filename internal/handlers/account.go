package handlers

import (
	"encoding/json"
	"net/http"
)

func AccountDetailsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user from context (set by JWT middleware)
		email := r.Context().Value("userEmail").(string)

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"id":   "user-" + email,
				"type": "users",
				"attributes": map[string]interface{}{
					"username":           email,
					"email":              email,
					"is-service-account": false,
					"avatar-url":         "https://www.gravatar.com/avatar/00000000000000000000000000000000?d=mp",
					"permissions": map[string]bool{
						"can-create-organizations": true,
						"can-change-email":         true,
						"can-change-username":      true,
					},
				},
				"links": map[string]string{
					"self": "/api/v2/users/user-" + email,
				},
			},
		}

		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(response)
	}
}
