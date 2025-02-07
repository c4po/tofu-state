package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type Session struct {
	ID        string
	UserEmail string
	ExpiresAt time.Time
}

var sessions = make(map[string]Session) // In-memory store, use Redis in production

func CreateSession(email string) string {
	b := make([]byte, 32)
	rand.Read(b)
	sessionID := base64.URLEncoding.EncodeToString(b)

	sessions[sessionID] = Session{
		ID:        sessionID,
		UserEmail: email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return sessionID
}

func GetSession(sessionID string) (Session, bool) {
	session, exists := sessions[sessionID]
	if !exists || time.Now().After(session.ExpiresAt) {
		return Session{}, false
	}
	return session, true
}
