package app_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetTxnsByTracker(t *testing.T) {
	testTrackerID := "test-tracker-id"
	emptyTransactions := []app.SharedTransaction{}
	sharedTransactions := []app.SharedTransaction{
		{
			ID: "test-shared-transaction",
		},
	}

	tests := map[string]struct {
		trackerID      string
		expectationsFn mock_app.MockAppFn
		wantTxns       []app.SharedTransaction
		wantCode       int
	}{
		"with an empty list of txns from the store, returns an empty list": {
			trackerID: testTrackerID,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTxnsByTracker(gomock.Eq(testTrackerID)).Return(emptyTransactions, nil).Times(1)
			},
			wantTxns: emptyTransactions,
			wantCode: http.StatusOK,
		},
		"with a list of txns from the store, returns the list": {
			trackerID: testTrackerID,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTxnsByTracker(gomock.Eq(testTrackerID)).Return(sharedTransactions, nil).Times(1)
			},
			wantTxns: sharedTransactions,
			wantCode: http.StatusOK,
		},
		"with a trackerID of a non-existent tracker, returns a 404": {
			trackerID: "non-existent-trackerID",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTxnsByTracker(gomock.Eq("non-existent-trackerID")).Return(nil, app.ErrDBItemNotFound).Times(1)
			},
			wantTxns: nil,
			wantCode: http.StatusNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			request := app.NewGetTxnsByTrackerRequest(tc.trackerID)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.GetTxnsByTracker)
			handler.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)
			if tc.wantTxns != nil {
				var got []app.SharedTransaction
				err := json.NewDecoder(response.Body).Decode(&got)
				if err != nil {
					t.Fatalf("error parsing response from server %q into slice of SharedTransactions, '%v'", response.Body, err)
				}
				assert.ElementsMatch(got, tc.wantTxns)
			}
		})
	}
}

func TestCreateSharedTxn(t *testing.T) {
	validSharedTxn := app.SharedTransaction{
		Participants: []string{"user1", "user2"},
		Amount:       123,
		Date:         123456,
		Shop:         "test-shop",
		Tracker:      "test-tracker-1",
		Category:     "test-category",
		Payer:        "user1",
	}

	tests := map[string]struct {
		transaction      app.SharedTransaction
		expectationsFn   mock_app.MockAppFn
		wantCode         int
		userInContext    string
		trackerInContext string
	}{
		"with a userID in the context that doesn't match one of the users in the transaction": {
			transaction: app.SharedTransaction{
				Participants: []string{"user1", "user2"},
				Amount:       123,
				Date:         123456,
				Shop:         "test-shop",
			},
			userInContext: "user-not-participating",
			wantCode:      http.StatusForbidden,
		},
		// TODO: handle a tracker not existing when trying to create a shared transaction
		"with a transaction missing a shop": {
			transaction: app.SharedTransaction{
				Participants: []string{"user1", "user2"},
				Amount:       123,
				Date:         123456,
			},
			userInContext: "user1",
			wantCode:      http.StatusBadRequest,
		},
		"with a transaction missing a date": {
			transaction: app.SharedTransaction{
				Participants: []string{"user1", "user2"},
				Amount:       123,
				Shop:         "test-shop",
			},
			userInContext: "user1",
			wantCode:      http.StatusBadRequest,
		},
		"with a transaction missing an amount": {
			transaction: app.SharedTransaction{
				Participants: []string{"user1", "user2"},
				Date:         123456,
				Shop:         "test-shop",
			},
			userInContext: "user1",
			wantCode:      http.StatusBadRequest,
		},
		"with a valid transaction and tracker in context": {
			transaction: validSharedTxn,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().CreateSharedTxn(validSharedTxn).Times(1)
			},
			userInContext:    "user1",
			trackerInContext: validSharedTxn.Tracker,
			wantCode:         http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			req := app.NewCreateSharedTxnRequest(tc.transaction)
			ctx := context.WithValue(req.Context(), app.CtxKeyUserID, tc.userInContext)
			ctx = context.WithValue(ctx, app.CtxKeyTrackerID, tc.trackerInContext)

			req = req.WithContext(ctx)
			response := httptest.NewRecorder()
			handler := http.HandlerFunc(a.CreateSharedTxn)
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}

func TestUpdateSharedTxn(t *testing.T) {
	validSharedTxn := app.SharedTransaction{
		Participants: []string{"user1", "user2"},
		Amount:       123,
		Date:         123456,
		Shop:         "test-shop",
		Tracker:      "test-tracker-1",
		Category:     "test-category",
		Payer:        "user1",
	}

	tests := map[string]struct {
		updatedTxn     app.SharedTransaction
		expectationsFn mock_app.MockAppFn
		wantCode       int
	}{
		"with a valid transaction, calls the store function successfully": {
			updatedTxn: validSharedTxn,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().UpdateSharedTxn(validSharedTxn).Times(1)
			},
			wantCode: http.StatusAccepted,
		},
		// when the store returns an error
		// with an invalid transaction
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			req := app.NewUpdateSharedTxnRequest(tc.updatedTxn)
			ctx := context.WithValue(req.Context(), app.CtxKeyUserID, tc.updatedTxn.Payer)
			ctx = context.WithValue(ctx, app.CtxKeyTrackerID, tc.updatedTxn.Tracker)
			ctx = context.WithValue(ctx, app.CtxKeyTransactionID, tc.updatedTxn.ID)
			req = req.WithContext(ctx)
			response := httptest.NewRecorder()
			handler := http.HandlerFunc(a.UpdateSharedTxn)
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}

func TestDeleteSharedTxn(t *testing.T) {
	testTxn := app.SharedTransaction{
		ID:           "test-txn-id",
		Participants: []string{"user1", "user2"},
		Tracker:      "test-tracker-1",
	}
	delTxnInput := app.DelSharedTxnInput{
		TxnID:        testTxn.ID,
		Participants: testTxn.Participants,
		Tracker:      testTxn.Tracker,
	}

	tests := map[string]struct {
		txn      app.SharedTransaction
		expectFn mock_app.MockAppFn
		wantCode int
	}{
		"calls the store function successfully": {
			txn: testTxn,
			expectFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().DeleteSharedTxn(delTxnInput).Times(1)
			},
			wantCode: http.StatusAccepted,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectFn)

			req := app.NewDeleteSharedTxnRequest(tc.txn)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.DeleteSharedTxn)
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}

func TestGetUnsettledTxnsByTracker(t *testing.T) {
	testTrackerID := "test-tracker-id"
	emptyTransactions := []app.SharedTransaction{}
	unsettledTxns := []app.SharedTransaction{{
		ID:        "test-unsettled-transaction",
		Unsettled: true,
	}}

	tests := map[string]struct {
		trackerID      string
		expectationsFn mock_app.MockAppFn
		wantTxns       []app.SharedTransaction
		wantCode       int
	}{
		"with an empty list of txns from the store, returns an empty list": {
			trackerID: testTrackerID,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetUnsettledTxnsByTracker(gomock.Eq(testTrackerID)).Return(emptyTransactions, nil).Times(1)
			},
			wantTxns: emptyTransactions,
			wantCode: http.StatusOK,
		},
		"with a list of txns from the store, returns the list": {
			trackerID: testTrackerID,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetUnsettledTxnsByTracker(gomock.Eq(testTrackerID)).Return(unsettledTxns, nil).Times(1)
			},
			wantTxns: unsettledTxns,
			wantCode: http.StatusOK,
		},
		"with a trackerID of a non-existent tracker, returns a 404": {
			trackerID: "non-existent-trackerID",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetUnsettledTxnsByTracker(gomock.Eq("non-existent-trackerID")).Return(nil, app.ErrDBItemNotFound).Times(1)
			},
			wantTxns: nil,
			wantCode: http.StatusNotFound,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			request := app.NewGetUnsettledTxnsByTrackerRequest(tc.trackerID)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.GetUnsettledTxnsByTracker)
			handler.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)
			if tc.wantTxns != nil {
				var got []app.SharedTransaction
				err := json.NewDecoder(response.Body).Decode(&got)
				if err != nil {
					t.Fatalf("error parsing response from server %q into slice of SharedTransactions, '%v'", response.Body, err)
				}
				assert.ElementsMatch(got, tc.wantTxns)
			}
		})
	}
}

func TestSettleTxns(t *testing.T) {
	testTransaction := app.SharedTransaction{
		ID:           "test-shared-txn-id",
		Participants: []string{"user-01", "user-02"},
		Tracker:      "test-tracker-id",
	}

	tests := map[string]struct {
		transactions   []app.SharedTransaction
		expectationsFn mock_app.MockAppFn
		wantCode       int
	}{
		"calls the store function successfully": {
			transactions: []app.SharedTransaction{testTransaction},
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().SettleTxns([]app.SharedTransaction{testTransaction}).Times(1)
			},
			wantCode: http.StatusAccepted,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			request := app.NewSettleTxnsRequest(t, tc.transactions)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.SettleTxns)
			handler.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}
