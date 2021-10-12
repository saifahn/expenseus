package sessions

import (
	"net/http"

	"github.com/gorilla/securecookie"
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

func (sm *SessionManager) StoreSession(rw http.ResponseWriter, r *http.Request) {

}
