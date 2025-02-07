package auth

import (
	"net/http"

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
