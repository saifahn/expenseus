package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyUser(t *testing.T) {
	t.Run("returns a 401 response when the user is not authorized", func(t *testing.T) {
		store := StubTransactionStore{}
		oauth := StubOauthConfig{}
		a := New(&store, &oauth, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetAllTransactionsRequest()
		response := httptest.NewRecorder()

		handler := a.VerifyUser(http.HandlerFunc(a.GetAllTransactions))
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("returns a 200 response when the user is authorized, and passes the request with the user ID in the context to the appropriate route", func(t *testing.T) {
		store := StubTransactionStore{transactions: map[string]Transaction{"1": TestSeanTransaction}}
		oauth := StubOauthConfig{}
		a := New(&store, &oauth, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetAllTransactionsRequest()
		// simulate a cookie session storage here
		request.AddCookie(&ValidCookie)
		response := httptest.NewRecorder()

		handler := a.VerifyUser(http.HandlerFunc(a.GetAllTransactions))
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var got []Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
		}

		assert.ElementsMatch(t, got, []Transaction{TestSeanTransaction})
	})
}

func TestLogOut(t *testing.T) {
	t.Run("session manager calls remove", func(t *testing.T) {
		store := StubTransactionStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		sessions := StubSessionManager{}
		frontend := "http://test.base"
		a := New(&store, &oauth, &sessions, frontend, &StubImageStore{})

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/logout", nil)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(a.LogOut)
		handler.ServeHTTP(response, request)

		assert.Equal(t, 1, sessions.removeCalls)
		assert.Equal(t, http.StatusTemporaryRedirect, response.Code)

		url, err := response.Result().Location()
		if err != nil {
			t.Fatalf("url couldn't be found: %v", err)
		}
		assert.Equal(t, frontend, url.String())
	})
}
