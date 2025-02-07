package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func AuthMiddleware(store sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, SessionName)
			if GetUserEmail(session) == "" {
				session.Values["return_to"] = r.URL.String()
				session.Save(r, w)
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func JWTMiddleware(secret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[JWT Middleware] Processing request for: %s\n", r.URL.Path)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				fmt.Println("[JWT Middleware] Missing Authorization header")
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			fmt.Println("[JWT Middleware] Found Authorization header, parsing token")
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				fmt.Printf("[JWT Middleware] Token validation failed: %v\n", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			fmt.Println("[JWT Middleware] Token is valid, checking claims")
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				fmt.Println("[JWT Middleware] Failed to parse token claims")
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			email, ok := claims["email"].(string)
			if !ok || email == "" {
				fmt.Println("[JWT Middleware] No valid email found in token claims")
				http.Error(w, "Invalid email in token", http.StatusUnauthorized)
				return
			}

			fmt.Printf("[JWT Middleware] Successfully authenticated user: %s\n", email)
			ctx := context.WithValue(r.Context(), "userEmail", email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
