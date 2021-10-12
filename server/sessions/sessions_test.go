package sessions

import (
	"net/http"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
)

var testHashKey = securecookie.GenerateRandomKey(64)
var testBlockKey = securecookie.GenerateRandomKey(32)

func TestValidateAuthorizedSession(t *testing.T) {
	sessions := New(testHashKey, testBlockKey)

	t.Run("returns false with no stored user id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)

		want := false
		got := sessions.ValidateAuthorizedSession(req)
		assert.Equal(t, got, want)
	})

	t.Run("returns false if the cookie value was encoded with a different secret", func(t *testing.T) {
		altHashKey := securecookie.GenerateRandomKey(64)
		altBlockKey := securecookie.GenerateRandomKey(32)
		s := securecookie.New(altHashKey, altBlockKey)
		encoded, err := s.Encode("expenseus-id", "test")
		if err != nil {
			t.Errorf("cookie could not be encoded: %v", err)
			return
		}

		cookie := &http.Cookie{
			Name:     "expenseus-id",
			Value:    encoded,
			Secure:   true,
			HttpOnly: true,
		}
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(cookie)

		want := false
		got := sessions.ValidateAuthorizedSession(req)
		assert.Equal(t, want, got)
	})

	t.Run("returns false if the cookie value is not encoded", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:     "expenseus-id",
			Value:    "test",
			Secure:   true,
			HttpOnly: true,
		}
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(cookie)

		want := false
		got := sessions.ValidateAuthorizedSession(req)
		assert.Equal(t, want, got)
	})

	t.Run("returns true with a stored user id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		encoded, err := sessions.cookies.Encode("expenseus-id", "test")
		if err != nil {
			t.Errorf("cookie could not be encoded: %v", err)
			return
		}

		cookie := &http.Cookie{
			Name:     "expenseus-id",
			Value:    encoded,
			Secure:   true,
			HttpOnly: true,
		}
		req.AddCookie(cookie)

		want := true
		got := sessions.ValidateAuthorizedSession(req)
		assert.Equal(t, want, got)
	})
}
