package sessions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/saifahn/expenseus/internal/app"
	"github.com/stretchr/testify/assert"
)

var (
	testHashKey      = securecookie.GenerateRandomKey(64)
	testBlockKey     = securecookie.GenerateRandomKey(32)
	testSessionValue = "testSession"
)

func TestValidate(t *testing.T) {
	sessions := New(testHashKey, testBlockKey)

	t.Run("returns false with no cookie", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)

		want := false
		got := sessions.Validate(req)
		assert.Equal(t, got, want)
	})

	t.Run("returns false if the cookie value was encoded with a different secret", func(t *testing.T) {
		altHashKey := securecookie.GenerateRandomKey(64)
		altBlockKey := securecookie.GenerateRandomKey(32)
		s := securecookie.New(altHashKey, altBlockKey)
		encoded, err := s.Encode(app.SessionCookieKey, testSessionValue)
		if err != nil {
			t.Errorf("cookie could not be encoded: %v", err)
			return
		}

		cookie := &http.Cookie{
			Name:     app.SessionCookieKey,
			Value:    encoded,
			Secure:   true,
			HttpOnly: true,
		}
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(cookie)

		want := false
		got := sessions.Validate(req)
		assert.Equal(t, want, got)
	})

	t.Run("returns false if the cookie value is not encoded", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:     app.SessionCookieKey,
			Value:    testSessionValue,
			Secure:   true,
			HttpOnly: true,
		}
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(cookie)

		want := false
		got := sessions.Validate(req)
		assert.Equal(t, want, got)
	})

	t.Run("returns true with a stored user id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		encoded, err := sessions.cookies.Encode(app.SessionCookieKey, testSessionValue)
		if err != nil {
			t.Errorf("cookie could not be encoded: %v", err)
			return
		}

		cookie := &http.Cookie{
			Name:     app.SessionCookieKey,
			Value:    encoded,
			Secure:   true,
			HttpOnly: true,
		}
		req.AddCookie(cookie)

		want := true
		got := sessions.Validate(req)
		assert.Equal(t, want, got)
	})
}

func TestSave(t *testing.T) {
	t.Run("given a request with a userid in context, stores the encoded id in a cookie of the appropriate name", func(t *testing.T) {
		sessions := New(testHashKey, testBlockKey)

		expectedUserID := testSessionValue
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), app.CtxKeyUserID, expectedUserID)
		req = req.WithContext(ctx)

		rw := httptest.NewRecorder()
		sessions.Save(rw, req)

		cookies := rw.Result().Cookies()
		for _, c := range cookies {
			if c.Name == app.SessionCookieKey {
				var userid string
				err := sessions.cookies.Decode(c.Name, c.Value, &userid)
				if err != nil {
					t.Fatalf("cookie could not be decoded: %v", err)
				}
				assert.Equal(t, expectedUserID, userid)
				return
			}
		}
		t.Fatalf("cookie with the expected name %q was not found", app.SessionCookieKey)
	})
}
