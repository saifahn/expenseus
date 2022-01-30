package sessions

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/saifahn/expenseus/internal/app"
)

type SessionManager struct {
	cookies securecookie.SecureCookie
}

func New(hashKey, blockKey []byte) *SessionManager {
	return &SessionManager{cookies: *securecookie.New(hashKey, blockKey)}
}

func (sm *SessionManager) Validate(r *http.Request) bool {
	cookie, err := r.Cookie(app.SessionCookieKey)
	if err != nil {
		return false
	}

	var userid string
	err = sm.cookies.Decode(app.SessionCookieKey, cookie.Value, &userid)
	return err == nil
}

func (sm *SessionManager) Save(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(app.CtxKeyUserID).(string)

	encoded, err := sm.cookies.Encode(app.SessionCookieKey, userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	cookie := &http.Cookie{
		Name:     app.SessionCookieKey,
		Value:    encoded,
		Secure:   true,
		HttpOnly: true,
		// one day
		MaxAge: 60 * 60 * 24,
	}

	http.SetCookie(rw, cookie)
}

func (sm *SessionManager) Remove(rw http.ResponseWriter, r *http.Request) {
	// overwrite the cookie with an expired cookie to delete it
	invalidCookie := &http.Cookie{
		Name:     app.SessionCookieKey,
		Value:    "deleted-cookie",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(-100),
	}

	http.SetCookie(rw, invalidCookie)
}

func (sm *SessionManager) GetUserID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(app.SessionCookieKey)
	if err != nil {
		return "", err
	}

	var userid string
	err = sm.cookies.Decode(app.SessionCookieKey, cookie.Value, &userid)
	if err != nil {
		return "", err
	}

	return userid, nil
}
