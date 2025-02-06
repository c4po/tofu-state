package handlers

import (
	"net/http"

	"github.com/c4po/tofu-state/internal/auth"
)

func HandleLogin(oidc *auth.OIDCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement login handler
	}
}

func HandleCallback(oidc *auth.OIDCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement callback handler
	}
}
