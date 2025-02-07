package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/c4po/tofu-state/internal/auth"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
)

func HandleLogin(oidc *auth.OIDCClient, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, auth.SessionName)

		// Generate random state
		state, err := generateRandomString(32)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Store state in session
		session.Values[auth.StateKey] = state
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Session save failed", http.StatusInternalServerError)
			return
		}

		// Build redirect URL
		redirectURL := fmt.Sprintf("%s://%s/callback",
			getProtocol(r),
			r.Host,
		)
		oidc.Config.RedirectURL = redirectURL

		http.Redirect(w, r, oidc.Config.AuthCodeURL(state), http.StatusFound)
	}
}

func getProtocol(r *http.Request) string {
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		return "https"
	}
	return "http"
}

func HandleCallback(oidc *auth.OIDCClient, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, auth.SessionName)

		// Verify state
		storedState, ok := session.Values[auth.StateKey].(string)
		if !ok || r.URL.Query().Get("state") != storedState {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		// Exchange code for token
		token, err := oidc.Config.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusUnauthorized)
			return
		}

		// Verify ID Token
		idToken, err := oidc.VerifyToken(r.Context(), token.Extra("id_token").(string))
		if err != nil {
			http.Error(w, "Failed to verify ID Token", http.StatusInternalServerError)
			return
		}

		// Extract claims
		var claims struct {
			Email string `json:"email"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Save user email in session
		session.Values["email"] = claims.Email
		session.Save(r, w)

		// Redirect to original URL
		returnTo, _ := session.Values["return_to"].(string)
		if returnTo == "" {
			returnTo = "/"
		}
		http.Redirect(w, r, returnTo, http.StatusFound)
	}
}

func HandleTokenRequest(oidc *auth.OIDCClient, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, auth.SessionName)
		email := auth.GetUserEmail(session)
		if email == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Generate and display token
		token, err := generateAPIToken(email)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<html>
				<body>
					<h1>API Token</h1>
					<pre>%s</pre>
					<p>Copy this token for use with OpenTofu</p>
				</body>
			</html>`, token)
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

func getEmailFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse token claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("no email found in token claims")
	}

	return email, nil
}
