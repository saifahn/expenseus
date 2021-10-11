package sessions

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testHashKey = "test-hash-key"
const testBlockKey = "test-block-key"

func TestValidateAuthorizedSession(t *testing.T) {
	t.Run("returns false with no stored user id", func(t *testing.T) {
		sessions := New([]byte(testHashKey), []byte(testBlockKey))
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)

		want := false
		got := sessions.ValidateAuthorizedSession(req)
		assert.Equal(t, got, want)
	})
}
