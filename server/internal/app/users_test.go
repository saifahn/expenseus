package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	store := StubTransactionStore{}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	user := TestSeanUser
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf(err.Error())
	}

	request := NewCreateUserRequest(userJSON)
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(app.CreateUser)
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusAccepted, response.Code)
	assert.Len(t, store.users, 1)
	assert.Contains(t, store.users, user)
}

func TestListUsers(t *testing.T) {
	store := StubTransactionStore{users: []User{TestSeanUser, TestTomomiUser}}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	request := NewGetAllUsersRequest()
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(app.ListUsers)
	handler.ServeHTTP(response, request)

	var got []User
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Len(t, got, len(store.users))
	assert.ElementsMatch(t, got, store.users)
}

func TestGetUserByID(t *testing.T) {
	t.Run("returns a users details if the user exists", func(t *testing.T) {
		store := StubTransactionStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		a := New(&store, &oauth, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetUserRequest(TestSeanUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(a.GetUser)
		handler.ServeHTTP(response, request)

		var got User
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, TestSeanUser, got)
	})
}

func TestGetSelf(t *testing.T) {
	t.Run("returns the user details from the stored session if the user exists", func(t *testing.T) {
		store := StubTransactionStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		a := New(&store, &oauth, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetSelfRequest()
		// add the user into the request cookie
		request.AddCookie(&ValidCookie)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(a.GetSelf)
		handler.ServeHTTP(response, request)

		// decode the response
		var got User
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, TestSeanUser, got)
	})

	t.Run("returns a 404 if the user does not exist", func(t *testing.T) {
		store := StubTransactionStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		a := New(&store, &oauth, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetSelfRequest()
		// add the user into the request cookie
		request.AddCookie(&http.Cookie{
			Name:  ValidCookie.Name,
			Value: "non-existent-user",
		})
		response := httptest.NewRecorder()
		handler := http.HandlerFunc(a.GetSelf)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.NotEqual(t, jsonContentType, response.Result().Header.Get("content-type"))
	})
}
