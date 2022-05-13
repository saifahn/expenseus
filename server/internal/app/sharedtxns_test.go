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

var testSharedTransaction = SharedTransaction{
	Participants: []string{"user1", "user2"},
	Shop:         "Test Shop",
	Amount:       123,
	Date:         123456,
}

func TestCreateSharedTxn(t *testing.T) {
	t.Run("CreateSharedTxn calls CreateSharedTxn with the transaction when passed a valid shared transaction", func(t *testing.T) {
		assert := assert.New(t)
		store := StubTransactionStore{}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		// make the request
		request := NewCreateSharedTxnRequest(testSharedTransaction, "user1")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateSharedTxn)
		handler.ServeHTTP(response, request)

		assert.Equal(http.StatusAccepted, response.Code)
		assert.Len(store.createSharedTxnCalls, 1)
		assert.Equal(testSharedTransaction, store.createSharedTxnCalls[0])
	})

	tests := map[string]struct {
		transaction   SharedTransaction
		wantCode      int
		userInContext string
	}{
		"with a userID in the context that doesn't match one of the users in the transaction": {
			transaction: SharedTransaction{
				Participants: []string{"user1", "user2"},
				Amount:       123,
				Date:         123456,
				Shop:         "test-shop",
			},
			userInContext: "user-not-participating",
			wantCode:      http.StatusForbidden,
		},
		"with a transaction missing a shop": {
			transaction: SharedTransaction{
				Participants: []string{"user1", "user2"},
				Amount:       123,
				Date:         123456,
			},
			userInContext: "user1",
			wantCode:      http.StatusBadRequest,
		},
		"with a transaction missing a date": {
			transaction: SharedTransaction{
				Participants: []string{"user1", "user2"},
				Amount:       123,
				Shop:         "test-shop",
			},
			userInContext: "user1",
			wantCode:      http.StatusBadRequest,
		},
		"with a transaction missing an amount": {
			transaction: SharedTransaction{
				Participants: []string{"user1", "user2"},
				Date:         123456,
				Shop:         "test-shop",
			},
			userInContext: "user1",
			wantCode:      http.StatusBadRequest,
		},
		"with a valid transaction": {
			transaction: SharedTransaction{
				Participants: []string{"user1", "user2"},
				Amount:       123,
				Date:         123456,
				Shop:         "test-shop",
			},
			userInContext: "user1",
			wantCode:      http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			store := StubTransactionStore{}
			app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

			request := NewCreateSharedTxnRequest(tc.transaction, tc.userInContext)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(app.CreateSharedTxn)
			handler.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}
