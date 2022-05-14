package app_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetTxnsByTracker(t *testing.T) {
	t.Run("get transactions by tracker ID calls the GetTxnsByTracker function", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_app.NewMockStore(ctrl)
		a := app.New(mockStore, &app.StubOauthConfig{}, &app.StubSessionManager{}, "", &app.StubImageStore{})

		trackerID := app.TestTracker.ID
		request := app.NewGetTxnsByTrackerRequest(t, trackerID)
		response := httptest.NewRecorder()

		mockStore.EXPECT().GetTxnsByTracker(gomock.Eq(trackerID)).Times(1)

		handler := http.HandlerFunc(a.GetTxnsByTracker)
		handler.ServeHTTP(response, request)
	})

	t.Run("returns a list of transactions when received from the store", func(t *testing.T) {
		assert := assert.New(t)
		ctrl := gomock.NewController(t)
		mockStore := mock_app.NewMockStore(ctrl)
		a := app.New(mockStore, &app.StubOauthConfig{}, &app.StubSessionManager{}, "", &app.StubImageStore{})

		trackerID := app.TestTracker.ID
		request := app.NewGetTxnsByTrackerRequest(t, trackerID)
		response := httptest.NewRecorder()

		wantTransactions := []app.SharedTransaction{}
		mockStore.EXPECT().GetTxnsByTracker(gomock.Eq(trackerID)).Return(wantTransactions, nil).Times(1)

		handler := http.HandlerFunc(a.GetTxnsByTracker)
		handler.ServeHTTP(response, request)

		var got []app.SharedTransaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of SharedTransactions, '%v'", response.Body, err)
		}

		// assert.Equal(app.jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(http.StatusOK, response.Code)
		assert.ElementsMatch(got, wantTransactions)
	})
	// TODO: I should add like custom mock values from the store and then do tests based on the expected behaviour
	// empty list of transactions
	// list of transactions
	// tracker doesn't exist
}

// var testSharedTransaction = app.SharedTransaction{
// 	Participants: []string{"user1", "user2"},
// 	Shop:         "Test Shop",
// 	Amount:       123,
// 	Date:         123456,
// }

// func TestCreateSharedTxn(t *testing.T) {
// 	t.Run("CreateSharedTxn calls the store's CreateSharedTxn with the transaction when passed a valid shared transaction", func(t *testing.T) {
// 		assert := assert.New(t)
// 		store := StubTransactionStore{}
// 		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

// 		request := NewCreateSharedTxnRequest(testSharedTransaction, "user1")
// 		response := httptest.NewRecorder()

// 		handler := http.HandlerFunc(app.CreateSharedTxn)
// 		handler.ServeHTTP(response, request)

// 		assert.Equal(http.StatusAccepted, response.Code)
// 		assert.Len(store.createSharedTxnCalls, 1)
// 		assert.Equal(testSharedTransaction, store.createSharedTxnCalls[0])
// 	})

// 	tests := map[string]struct {
// 		transaction   SharedTransaction
// 		wantCode      int
// 		userInContext string
// 	}{
// 		"with a userID in the context that doesn't match one of the users in the transaction": {
// 			transaction: SharedTransaction{
// 				Participants: []string{"user1", "user2"},
// 				Amount:       123,
// 				Date:         123456,
// 				Shop:         "test-shop",
// 			},
// 			userInContext: "user-not-participating",
// 			wantCode:      http.StatusForbidden,
// 		},
// 		"with a transaction missing a shop": {
// 			transaction: SharedTransaction{
// 				Participants: []string{"user1", "user2"},
// 				Amount:       123,
// 				Date:         123456,
// 			},
// 			userInContext: "user1",
// 			wantCode:      http.StatusBadRequest,
// 		},
// 		"with a transaction missing a date": {
// 			transaction: SharedTransaction{
// 				Participants: []string{"user1", "user2"},
// 				Amount:       123,
// 				Shop:         "test-shop",
// 			},
// 			userInContext: "user1",
// 			wantCode:      http.StatusBadRequest,
// 		},
// 		"with a transaction missing an amount": {
// 			transaction: SharedTransaction{
// 				Participants: []string{"user1", "user2"},
// 				Date:         123456,
// 				Shop:         "test-shop",
// 			},
// 			userInContext: "user1",
// 			wantCode:      http.StatusBadRequest,
// 		},
// 		"with a valid transaction": {
// 			transaction: SharedTransaction{
// 				Participants: []string{"user1", "user2"},
// 				Amount:       123,
// 				Date:         123456,
// 				Shop:         "test-shop",
// 			},
// 			userInContext: "user1",
// 			wantCode:      http.StatusAccepted,
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			assert := assert.New(t)
// 			store := StubTransactionStore{}
// 			app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

// 			request := NewCreateSharedTxnRequest(tc.transaction, tc.userInContext)
// 			response := httptest.NewRecorder()

// 			handler := http.HandlerFunc(app.CreateSharedTxn)
// 			handler.ServeHTTP(response, request)

// 			assert.Equal(tc.wantCode, response.Code)
// 		})
// 	}
// }

// func TestGetUnsettledTxnsByTracker(t *testing.T) {
// 	t.Run("GetUnsettledTxnsByTracker calls the store's GetUnsettledTxnsByTracker with a given trackerID", func(t *testing.T) {
// 		assert := assert.New(t)
// 		store := StubTransactionStore{}
// 		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})
// 		wantTrackerID := "testTracker"

// 		request := NewGetUnsettledTxnsByTrackerRequest(wantTrackerID)
// 		response := httptest.NewRecorder()

// 		handler := http.HandlerFunc(app.GetUnsettledTxnsByTracker)
// 		handler.ServeHTTP(response, request)

// 		assert.Equal(http.StatusOK, response.Code)
// 		assert.Len(store.getUnsettledTxnsByTrackerCalls, 1)
// 		assert.Equal(wantTrackerID, store.getUnsettledTxnsByTrackerCalls[0])
// 	})
// }
