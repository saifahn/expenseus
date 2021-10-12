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

	// TODO: no user id separate from no cookie?
	// TODO: error with different hash key decoding?

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
