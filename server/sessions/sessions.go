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
	// sm.cookies.Decode()
	return false
}

func (sm *SessionManager) StoreSession(rw http.ResponseWriter, r *http.Request) {

}
