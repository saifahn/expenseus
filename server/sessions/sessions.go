package sessions

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/saifahn/expenseus"
)

type SessionManager struct {
	cookies securecookie.SecureCookie
}

func New(hashKey, blockKey []byte) *SessionManager {
	return &SessionManager{cookies: *securecookie.New(hashKey, blockKey)}
}

func (sm *SessionManager) ValidateAuthorizedSession(r *http.Request) bool {
	cookie, err := r.Cookie("expenseus-id")
	if err != nil {
		return false
	}

	var userid string
	err = sm.cookies.Decode("expenseus-id", cookie.Value, &userid)
	return err == nil
}

func (sm *SessionManager) SaveSession(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(expenseus.CtxKeyUserID).(string)

	encoded, err := sm.cookies.Encode("expenseus-id", userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	cookie := &http.Cookie{
		Name:     "expenseus-id",
		Value:    encoded,
		Secure:   true,
		HttpOnly: true,
		// one day
		MaxAge: 60 * 60 * 24,
	}

	http.SetCookie(rw, cookie)
}
