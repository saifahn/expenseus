package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifahn/expenseus/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestCreatingUsersAndRetrievingThem(t *testing.T) {
	t.Run("a valid cookie must be provided in order to create a user, GetSelf will read the cookie and attempt to get the user from the ID within, and a user can be retrieved by ID", func(t *testing.T) {
		router, tearDownDB := SetUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// TRY to create a user WITHOUT a valid cookie
		userJSON, err := json.Marshal(TestSeanUser)
		if err != nil {
			t.Fatalf("failed to marshal the user JSON: %v", err)
		}
		response := httptest.NewRecorder()
		request := app.NewCreateUserRequest(userJSON)
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusUnauthorized, response.Code)

		// use a VALID cookie
		CreateUser(t, TestSeanUser, router)

		// TRY GetSelf with different ID in the cookie
		// should not work as the userID from the cookie does not exist
		response = httptest.NewRecorder()
		request = app.NewGetSelfRequest()
		request.AddCookie(CreateCookie("not-real-id"))

		router.ServeHTTP(response, request)
		assert.Equal(http.StatusNotFound, response.Code)

		// use a cookie with the SAME ID
		response = httptest.NewRecorder()
		request = app.NewGetSelfRequest()
		request.AddCookie(CreateCookie(TestSeanUser.ID))
		router.ServeHTTP(response, request)

		var userGot app.User
		err = json.NewDecoder(response.Body).Decode(&userGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(TestSeanUser, userGot)

		// GET the specifically created user from the db by ID
		response = httptest.NewRecorder()
		request = app.NewGetUserRequest(TestSeanUser.ID)
		request.AddCookie(CreateCookie(TestSeanUser.ID))
		router.ServeHTTP(response, request)

		err = json.NewDecoder(response.Body).Decode(&userGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(TestSeanUser, userGot)
	})

	t.Run("multiple users can be created and retrieved with a request to the GetAllUsers route", func(t *testing.T) {
		router, tearDownDB := SetUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create TWO users
		CreateUser(t, TestSeanUser, router)
		CreateUser(t, TestTomomiUser, router)

		// GET all users
		response := httptest.NewRecorder()
		request := app.NewGetAllUsersRequest()
		request.AddCookie(CreateCookie(TestSeanUser.ID))
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		// ensure that they contain the two users
		var usersGot []app.User
		err := json.NewDecoder(response.Body).Decode(&usersGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Users: %v", response.Body, err)
		}
		assert.ElementsMatch(usersGot, []app.User{TestSeanUser, TestTomomiUser})
	})
}

func createTestTransaction(t *testing.T, r http.Handler, td app.Transaction, userid string) {
	payload := app.MakeTxnRequestPayload(td)
	request := app.NewCreateTransactionRequest(payload)
	request.AddCookie(CreateCookie(userid))
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	assert.Equal(t, http.StatusAccepted, response.Code)
}
