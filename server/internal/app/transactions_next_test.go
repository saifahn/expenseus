package app_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

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

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/transactions/%s", tc.transactionId), nil)
			assert.NoError(err)
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
