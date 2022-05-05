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
	Shop:   "Test Shop",
	Amount: 123,
	Date:   123456,
}

func TestCreateSharedTxn(t *testing.T) {
	t.Run("CreateSharedTxn calls CreateSharedTxn with the transaction when passed a valid shared transaction", func(t *testing.T) {
		assert := assert.New(t)
		store := StubTransactionStore{}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		// make the request
		request := NewCreateSharedTxnRequest(testSharedTransaction)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateSharedTxn)
		handler.ServeHTTP(response, request)

		assert.Equal(http.StatusAccepted, response.Code)
		assert.Len(store.createSharedTxnCalls, 1)
		assert.Equal(testSharedTransaction, store.createSharedTxnCalls[0])
	})

	// tests := map[string]struct {
	// 	transaction SharedTransaction
	// 	wantCode    int
	// }{
	// 	"with a valid transaction": {
	// 		transaction: testSharedTransaction,
	// 		wantCode:    http.StatusOK,
	// 	},
	// }

	// for name, tc := range tests {
	// 	t.Run(name, func(t *testing.T) {
	// 		request := NewCreateSharedTxnRequest(t, tc.transaction)
	// 	})
	// }

	// t.Run("CreateSharedTxn calls the CreateSharedTxn function", func(t *testing.T) {
	// shared transaction
	// })
	// table test
	// with no user in cookie, error
	// with user in the cookie but wrong user, error
	// with user in the cookie correct, correct
	// with an invalid transaction
	// with a valid transaction
}
