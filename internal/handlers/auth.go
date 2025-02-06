package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/c4po/tofu-state/internal/auth"
)

func HandleLogin(oidc *auth.OIDCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate random state
		state, err := generateRandomString(32)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Store state in cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   true,
		})

		// Redirect to OIDC provider
		http.Redirect(w, r, oidc.Config.AuthCodeURL(state), http.StatusFound)
	}
}

func HandleCallback(oidc *auth.OIDCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify state
		stateCookie, err := r.Cookie("oauth_state")
		if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		// Exchange code for token
		token, err := oidc.Config.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusUnauthorized)
			return
		}

		// Create session or return token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": token.AccessToken,
			"token_type":   token.TokenType,
		})
	}
}

func HandleToken(oidc *auth.OIDCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement OAuth2 token endpoint
		// This should handle the token exchange using your OIDC client
		// Example:
		token, err := oidc.Config.Exchange(r.Context(), r.FormValue("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(token)
	}
}

func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return base64.URLEncoding.EncodeToString(b), err
}
