package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	testTxnID := "test-txn-id"
	tests := map[string]struct {
		txnID          string
		expectationsFn mock_app.MockAppFn
		wantTxns       app.Transaction
		wantCode       int
	}{
		"calls the store function to get the transaction": {
			txnID: testTxnID,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransaction(gomock.Eq(testTxnID)).Return(app.Transaction{ID: testTxnID}, nil).Times(1)
			},
			wantTxns: app.Transaction{ID: testTxnID},
			wantCode: http.StatusOK,
		},
		// "returns a response without an image":                  {},
		// "returns an image url if the transaction has an image": {},
		// "returns 404 on a non-existent transaction":            {},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			req := app.NewGetTransactionRequest(tc.txnID)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.GetTransaction)
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}

func TestDeleteTransaction(t *testing.T) {
	// testTransaction := app.Transaction{}

	tests := map[string]struct {
		transactionId  string
		user           string
		expectationsFn mock_app.MockAppFn
		wantCode       int
	}{
		"calls the store function to delete the transaction": {
			transactionId: "test-transaction-id",
			user:          "test-user",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().DeleteTransaction(gomock.Eq("test-transaction-id"), gomock.Eq("test-user")).Return(nil).Times(1)
			},
			wantCode: http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			req := app.NewDeleteTransactionRequest(tc.transactionId)
			ctx := context.WithValue(req.Context(), app.CtxKeyUserID, tc.user)
			ctx = context.WithValue(ctx, app.CtxKeyTransactionID, tc.transactionId)

			req = req.WithContext(ctx)
			response := httptest.NewRecorder()
			handler := http.HandlerFunc(a.DeleteTransaction)
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}
