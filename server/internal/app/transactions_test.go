package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllTransactions(t *testing.T) {
	t.Run("gets all transactions with one transaction", func(t *testing.T) {
		wantedTransactions := []Transaction{
			TestTomomiTransaction,
		}
		store := StubTransactionStore{
			users: []User{},
			transactions: map[string]Transaction{
				"9281": TestTomomiTransaction,
			},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetAllTransactionsRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetAllTransactions)
		handler.ServeHTTP(response, request)

		var got []Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(wantedTransactions), len(got))
		assert.ElementsMatch(t, got, wantedTransactions)
	})

	t.Run("gets all transactions with more than one transaction", func(t *testing.T) {
		wantedTransactions := []Transaction{
			TestSeanTransaction, TestTomomiTransaction, TestTomomiTransaction2,
		}
		store := StubTransactionStore{
			users: []User{},
			transactions: map[string]Transaction{
				"1":     TestSeanTransaction,
				"9281":  TestTomomiTransaction,
				"14928": TestTomomiTransaction2,
			},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetAllTransactionsRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetAllTransactions)
		handler.ServeHTTP(response, request)

		var got []Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(wantedTransactions), len(got))
		assert.ElementsMatch(t, got, wantedTransactions)
	})
}
