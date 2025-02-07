package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/c4po/tofu-state/internal/auth"
	"github.com/c4po/tofu-state/internal/config"
	"github.com/golang-jwt/jwt"
)

func HandleLogin(oidc *auth.OIDCClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Build redirect URL from request
		redirectURL := fmt.Sprintf("%s://%s%s",
			getProtocol(r),
			r.Host,
			oidc.Config.RedirectURL,
		)
		oidc.Config.RedirectURL = redirectURL

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

func getProtocol(r *http.Request) string {
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		return "https"
	}
	return "http"
}

func HandleCallback(oidc *auth.OIDCClient, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify state
		stateCookie, err := r.Cookie("oauth_state")
		if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		// Check if this is a token request
		isTokenRequest := strings.HasPrefix(stateCookie.Value, "token_request:")

		// Exchange code for token
		token, err := oidc.Config.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusUnauthorized)
			return
		}

		if isTokenRequest {
			// Generate API token
			apiToken, err := generateAPIToken(token.Extra("id_token").(string))
			if err != nil {
				http.Error(w, "Failed to generate API token", http.StatusInternalServerError)
				return
			}

			// Display token directly in browser
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `
				<html><body>
					<h1>OpenTofu API Token</h1>
					<pre style="background:#eee;padding:1rem">%s</pre>
					<p>Copy this token to use with OpenTofu:</p>
					<code>tofu login %s -token="%s"</code>
				</body></html>`,
				apiToken, cfg.ExternalHost, apiToken)
			return
		}

		// Get user info from ID token
		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No ID token", http.StatusInternalServerError)
			return
		}

		idToken, err := oidc.VerifyToken(r.Context(), rawIDToken)
		if err != nil {
			http.Error(w, "Invalid ID token", http.StatusUnauthorized)
			return
		}

		var claims struct {
			Email string `json:"email"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Create session
		sessionID := auth.CreateSession(claims.Email)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionID,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
		})

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

func HandleTokenRequest(oidc *auth.OIDCClient, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check existing session
		sessionCookie, err := r.Cookie("session_token")
		if err == nil {
			if session, valid := auth.GetSession(sessionCookie.Value); valid {
				// Generate token for logged-in user
				token, err := generateAPIToken(session.UserEmail)
				if err == nil {
					w.Header().Set("Content-Type", "text/html")
					fmt.Fprintf(w, `<html><body>
						<h1>OpenTofu API Token</h1>
						<pre style="background:#eee;padding:1rem">%s</pre>
						<p>Copy this token to use with OpenTofu:</p>
						<code>tofu login %s -token="%s"</code>
					</body></html>`, token, cfg.ExternalHost, token)
					return
				}
			}
		}

		// No valid session - start OIDC flow
		state, err := generateRandomString(16)
		if err != nil {
			http.Error(w, "Failed to generate state", http.StatusInternalServerError)
			return
		}
		state = fmt.Sprintf("token_request:%s", state)

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

func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateAPIToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   email,
		"email": email,
		"exp":   time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	// Store token in database
	// err := storageBackend.StoreToken(email, claims)
	// if err != nil {
	// 	return "", err
	// }

	// Sign token with server secret
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte("your-secret-key")) // Use proper secret management
}
