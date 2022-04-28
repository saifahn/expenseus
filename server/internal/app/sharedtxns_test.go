package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTransactionsByTracker(t *testing.T) {
	assert := assert.New(t)
	store := StubTransactionStore{}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	t.Run("get transactions by tracker ID calls the GetTxnsByTracker function", func(t *testing.T) {
		trackerID := TestTracker.ID
		request := NewGetTxnsByTrackerRequest(t, trackerID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTxnsByTracker)
		handler.ServeHTTP(response, request)

		assert.Len(store.getTxnsByTrackerCalls, 1)
		assert.Equal(trackerID, store.getTxnsByTrackerCalls[0])
	})
	// TODO: I should add like custom mock values from the store and then do tests based on the expected behaviour
}
