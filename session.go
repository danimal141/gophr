package main

import (
	"net/http"
	"net/url"
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

func RequestSession(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}

	session, err := globalSessionStore.Find(cookie.Value)
	if err != nil {
		panic(err)
	}
	if session == nil {
		return nil
	}
	if session.Expired() {
		globalSessionStore.Delete(session)
		return nil
	}
	return session
}

func RequestUser(r *http.Request) *User {
	session := RequestSession(r)
	if session == nil || session.UserID == "" {
		return nil
	}

	user, err := globalUserStore.Find(session.UserID)
	if err != nil {
		panic(err)
	}
	return user
}

func RequireLogin(w http.ResponseWriter, r *http.Request) {
	if RequestUser(r) != nil {
		return
	}
	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))
	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
}

func (session *Session) Expired() bool {
	return session.Expiry.Before(time.Now())
}
