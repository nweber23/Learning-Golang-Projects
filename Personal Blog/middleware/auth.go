package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	ID        string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

var store = &SessionStore{
	sessions: make(map[string]*Session),
}

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateSession(username string) (string, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return "", err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	session := &Session{
		ID:        sessionID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	store.sessions[sessionID] = session
	return sessionID, nil
}

func GetSession(id string) (*Session, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	session, exists := store.sessions[id]
	if !exists {
		return nil, errors.New("session not found")
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}
	return session, nil
}

func DeleteSession(id string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.sessions, id)
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		session, err := GetSession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		}
		next(w, r)
	}
}

func SetSessionCookie(w http.ResponseWriter, id string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   id,
		Path: 	"/",
		HttpOnly: true,
		MaxAge:   86400,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
	})
}