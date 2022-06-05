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

func TestGetTransaction(t *testing.T) {
	testTxnID := "test-txn-id"

	txnWithImageKey := app.Transaction{
		ID: testTxnID,
		TransactionDetails: app.TransactionDetails{
			ImageKey: "test-image-key",
		},
	}

	txnWithImageURL := app.Transaction{
		ID:                 txnWithImageKey.ID,
		ImageURL:           "test-image-url",
		TransactionDetails: txnWithImageKey.TransactionDetails,
	}

	tests := map[string]struct {
		txnID          string
		expectationsFn mock_app.MockAppFn
		wantTxn        *app.Transaction
		wantCode       int
	}{
		"calls the store function to get the transaction": {
			txnID: testTxnID,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransaction(gomock.Eq(testTxnID)).Return(app.Transaction{ID: testTxnID}, nil).Times(1)
			},
			wantTxn:  &app.Transaction{ID: testTxnID},
			wantCode: http.StatusOK,
		},
		"with a txnID of a transaction with an image": {
			txnID: testTxnID,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransaction(gomock.Eq(testTxnID)).Return(txnWithImageKey, nil).Times(1)
				ma.MockImages.EXPECT().AddImageToTransaction(gomock.Eq(txnWithImageKey)).Return(txnWithImageURL, nil)
			},
			wantTxn:  &txnWithImageURL,
			wantCode: http.StatusOK,
		},
		"returns 404 on a non-existent transaction": {
			txnID: "non-existent-txn",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransaction(gomock.Eq("non-existent-txn")).Return(app.Transaction{}, app.ErrDBItemNotFound).Times(1)
			},
			wantTxn:  nil,
			wantCode: http.StatusNotFound,
		},
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

			if tc.wantTxn != nil {
				var got app.Transaction
				err := json.NewDecoder(response.Body).Decode(&got)
				assert.NoError(err)
				assert.Equal(*tc.wantTxn, got)
			}
		})
	}
}

func TestGetTxnsByUser(t *testing.T) {
	testTxn := app.Transaction{ID: "txn-01"}
	testTxnWithImageKey := app.Transaction{
		ID:                 "txn-image-key",
		TransactionDetails: app.TransactionDetails{ImageKey: "test-image-key"},
	}
	testTxnWithImageURL := app.Transaction{
		ID:                 testTxnWithImageKey.ID,
		ImageURL:           "test-image-url",
		TransactionDetails: testTxnWithImageKey.TransactionDetails,
	}

	tests := map[string]struct {
		user           string
		expectationsFn mock_app.MockAppFn
		wantTxns       []app.Transaction
		wantCode       int
	}{
		"with a user that has one transaction": {
			user: "one-txn-user",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransactionsByUser(gomock.Eq("one-txn-user")).Return([]app.Transaction{testTxn}, nil).Times(1)
			},
			wantTxns: []app.Transaction{testTxn},
			wantCode: http.StatusOK,
		},
		"with a user that has 0 transactions": {
			user: "zero-txn-user",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransactionsByUser(gomock.Eq("zero-txn-user")).Return([]app.Transaction{}, nil).Times(1)
			},
			wantTxns: []app.Transaction{},
			wantCode: http.StatusOK,
		},
		"with transactions that have images": {
			user: "txn-with-image-user",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransactionsByUser(gomock.Eq("txn-with-image-user")).Return([]app.Transaction{testTxnWithImageKey}, nil).Times(1)
				ma.MockImages.EXPECT().AddImageToTransaction(gomock.Eq(testTxnWithImageKey)).Return(testTxnWithImageURL, nil).Times(1)
			},
			wantTxns: []app.Transaction{testTxnWithImageURL},
			wantCode: http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			req := app.NewGetTransactionsByUserRequest(tc.user)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.GetTransactionsByUser)
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)

			var got []app.Transaction
			err := json.NewDecoder(response.Body).Decode(&got)
			assert.NoError(err)
			assert.Equal(tc.wantTxns, got)
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
