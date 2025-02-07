package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SessionName = "tofu_state_session"
	UserKey     = "user_email"
	StateKey    = "oauth_state"
)

func GetUserEmail(s *sessions.Session) string {
	if val, ok := s.Values[UserKey].(string); ok {
		return val
	}
	return ""
}

func SetUserSession(w http.ResponseWriter, r *http.Request, store sessions.Store, email string) error {
	session, _ := store.Get(r, SessionName)
	session.Values[UserKey] = email
	return session.Save(r, w)
}

func ClearSession(w http.ResponseWriter, r *http.Request, store sessions.Store) error {
	session, _ := store.Get(r, SessionName)
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
