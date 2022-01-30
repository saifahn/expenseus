package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOauthLogin(t *testing.T) {
	store := StubTransactionStore{}
	oauth := StubOauthConfig{}
	app := New(&store, &oauth, &StubSessionManager{}, "", &StubImageStore{})

	request, err := http.NewRequest(http.MethodGet, "/api/v1/login_google", nil)
	if err != nil {
		t.Fatalf("request could not be created, %v", err)
	}
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(app.OauthLogin)
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusTemporaryRedirect, response.Code)
	// are these even good assertions to have?
	expectedURL := fmt.Sprintf("/api/v1/%s", oauthProviderMockURL)
	assert.Equal(t, expectedURL, response.Header().Get("Location"))
	// assert AuthCodeURL was called
	assert.Len(t, oauth.AuthCodeURLCalls, 1)
}

func TestOauthCallback(t *testing.T) {
	t.Run("creates a user when user doesn't exist yet and creates a new session with the user", func(t *testing.T) {
		store := StubTransactionStore{users: []User{}}
		oauth := StubOauthConfig{}
		sessions := StubSessionManager{}
		frontend := "http://a.test"
		app := New(&store, &oauth, &sessions, frontend, &StubImageStore{})

		request := NewGoogleCallbackRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.OauthCallback)
		handler.ServeHTTP(response, request)

		// expect a new user to be added to the store, GetInfoAndGenerateUser has been stubbed to generate TestSeanUser
		expected := []User{TestSeanUser}
		assert.Len(t, store.users, 1)
		assert.ElementsMatch(t, expected, store.users)

		assert.Len(t, sessions.saveCalls, 1)
		assert.Equal(t, sessions.saveCalls[0], TestSeanUser.ID)

		// get routed to the base page for now
		url, err := response.Result().Location()
		if err != nil {
			t.Fatalf("url couldn't be found: %v", err)
		}
		assert.Equal(t, frontend, url.String())
	})

	t.Run("doesn't create a new user when the user already exists, and saves the session with the user in the context", func(t *testing.T) {
		store := StubTransactionStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		sessions := StubSessionManager{}
		frontend := "http://another.test"
		app := New(&store, &oauth, &sessions, frontend, &StubImageStore{})

		request := NewGoogleCallbackRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.OauthCallback)
		handler.ServeHTTP(response, request)

		expected := []User{TestSeanUser}
		assert.Len(t, store.users, 1)
		assert.ElementsMatch(t, expected, store.users)

		assert.Len(t, sessions.saveCalls, 1)
		// the callback will add a context of the appropriate user id
		assert.Equal(t, sessions.saveCalls[0], TestSeanUser.ID)

		// expect to get routed to the main welcome page
		url, err := response.Result().Location()
		if err != nil {
			t.Fatalf("url couldn't be found: %v", err)
		}
		assert.Equal(t, frontend, url.String())
	})
}

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
