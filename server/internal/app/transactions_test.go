package app_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	testTxnID := "test-txn-id"
	testUserID := "test-user-id"

	txnWithImageKey := app.Transaction{
		ID:       testTxnID,
		ImageKey: "test-image-key",
	}

	txnWithImageURL := app.Transaction{
		ID:       txnWithImageKey.ID,
		ImageURL: "test-image-url",
		ImageKey: txnWithImageKey.ImageKey,
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
				ma.MockStore.EXPECT().GetTransaction(testUserID, gomock.Eq(testTxnID)).Return(app.Transaction{ID: testTxnID}, nil).Times(1)
			},
			wantTxn:  &app.Transaction{ID: testTxnID},
			wantCode: http.StatusOK,
		},
		"with a txnID of a transaction with an image": {
			txnID: testTxnID,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransaction(testUserID, gomock.Eq(testTxnID)).Return(txnWithImageKey, nil).Times(1)
				ma.MockImages.EXPECT().AddImageToTransaction(gomock.Eq(txnWithImageKey)).Return(txnWithImageURL, nil)
			},
			wantTxn:  &txnWithImageURL,
			wantCode: http.StatusOK,
		},
		"returns 404 on a non-existent transaction": {
			txnID: "non-existent-txn",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetTransaction(testUserID, gomock.Eq("non-existent-txn")).Return(app.Transaction{}, app.ErrDBItemNotFound).Times(1)
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
			ctx := context.WithValue(req.Context(), app.CtxKeyUserID, testUserID)
			ctx = context.WithValue(ctx, app.CtxKeyTransactionID, tc.txnID)
			req = req.WithContext(ctx)
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
		ID:       "txn-image-key",
		ImageKey: "test-image-key",
	}
	testTxnWithImageURL := app.Transaction{
		ID:       testTxnWithImageKey.ID,
		ImageURL: "test-image-url",
		ImageKey: testTxnWithImageKey.ImageKey,
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

func TestCreateTransaction(t *testing.T) {
	testTxnDetails := app.Transaction{
		Name:     "test-txn",
		Amount:   123,
		Date:     123456,
		Category: "test.category",
	}

	testImgTxnDetails := app.Transaction{
		Name:     "test-txn",
		Amount:   123,
		Date:     123456,
		ImageKey: "test-image-key",
		Category: "test.category",
	}

	tests := map[string]struct {
		txnDetails     app.Transaction
		user           string
		expectationsFn mock_app.MockAppFn
		wantCode       int
		withImg        bool
	}{
		"with a valid transaction without an image": {
			txnDetails: testTxnDetails,
			user:       "user-01",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().CreateTransaction(gomock.Eq(app.Transaction{
					Name:     testTxnDetails.Name,
					Amount:   testTxnDetails.Amount,
					Date:     testTxnDetails.Date,
					Category: testTxnDetails.Category,
					UserID:   "user-01",
				})).Return(nil).Times(1)
			},
			wantCode: http.StatusAccepted,
		},
		"with a transaction with an image that fails the image check": {
			txnDetails: testImgTxnDetails,
			user:       "user-02",
			expectationsFn: func(ma *mock_app.App) {
				// can't assert the file exactly, so use any
				ma.MockImages.EXPECT().Validate(gomock.Any()).Return(false, nil).Times(1)
			},
			wantCode: http.StatusUnprocessableEntity,
			withImg:  true,
		},
		"with a transaction with an image that passes the image check and successfully uploads": {
			txnDetails: testImgTxnDetails,
			user:       "user-02",
			expectationsFn: func(ma *mock_app.App) {
				ma.MockImages.EXPECT().Validate(gomock.Any()).Return(true, nil).Times(1)
				ma.MockImages.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(testImgTxnDetails.ImageKey, nil)
				ma.MockStore.EXPECT().CreateTransaction(gomock.Eq(app.Transaction{
					Name:     testImgTxnDetails.Name,
					Amount:   testImgTxnDetails.Amount,
					Date:     testImgTxnDetails.Date,
					ImageKey: testImgTxnDetails.ImageKey,
					Category: testImgTxnDetails.Category,
					UserID:   "user-02",
				})).Times(1)
			},
			withImg:  true,
			wantCode: http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			var payload map[string]io.Reader
			if tc.withImg {
				testFile, err := os.CreateTemp("", "example-file")
				if err != nil {
					t.Fatal(err)
				}
				defer testFile.Close()
				defer os.Remove(testFile.Name())

				payload = map[string]io.Reader{
					"transactionName": strings.NewReader(testImgTxnDetails.Name),
					"amount":          strings.NewReader("123"),
					"date":            strings.NewReader("123456"),
					"category":        strings.NewReader(testImgTxnDetails.Category),
					"image":           testFile,
				}
			} else {
				payload = app.MakeTxnRequestPayload(tc.txnDetails)
			}

			req := app.NewCreateTransactionRequest(payload)
			req = req.WithContext(context.WithValue(req.Context(), app.CtxKeyUserID, tc.user))
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.CreateTransaction)
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}

func TestDeleteTransaction(t *testing.T) {
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

func TestUpdateTransaction(t *testing.T) {
	updateTxnInput := app.Transaction{
		Name:     "test-transaction-name",
		Amount:   123,
		Date:     123456,
		Category: "test.category",
	}

	tests := map[string]struct {
		user           string
		txnID          string
		txnDetails     app.Transaction
		expectationsFn mock_app.MockAppFn
		wantCode       int
	}{
		"calls the store function to update the transaction": {
			user:       "test-user",
			txnID:      "test-transaction-id",
			txnDetails: updateTxnInput,
			expectationsFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().UpdateTransaction(app.Transaction{
					ID:       "test-transaction-id",
					UserID:   "test-user",
					Name:     updateTxnInput.Name,
					Amount:   updateTxnInput.Amount,
					Date:     updateTxnInput.Date,
					Category: updateTxnInput.Category,
				}).Return(nil).Times(1)
			},
			wantCode: http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectationsFn)

			req := app.NewUpdateTransactionRequest(tc.txnDetails)
			ctx := context.WithValue(req.Context(), app.CtxKeyUserID, tc.user)
			ctx = context.WithValue(ctx, app.CtxKeyTransactionID, tc.txnID)
			req = req.WithContext(ctx)

			response := httptest.NewRecorder()
			handler := http.HandlerFunc(a.UpdateTransaction)
			handler.ServeHTTP(response, req)
			assert.Equal(tc.wantCode, response.Code)
		})
	}
}
