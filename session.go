package main

import (
	"net/http"
	"time"
)

type Session struct {
	ID     string
	UserID string
	Expiry time.Time
}

const (
	sessionLifeTime   = 24 * 3 * time.Hour
	sessionCookieName = "GophrSession"
	sessionIDLength   = 20
)

func NewSession(w http.ResponseWriter) (*Session, error) {
	id, err := GenerateID("sess", sessionIDLength)
	if err != nil {
		return nil, err
	}

	expiry := time.Now().Add(sessionLifeTime)
	session := &Session{
		ID:     id,
		Expiry: expiry,
	}
	cookie := &http.Cookie{
		Name:    sessionCookieName,
		Value:   session.ID,
		Expires: expiry,
	}
	http.SetCookie(w, cookie)
	return session, nil
}
