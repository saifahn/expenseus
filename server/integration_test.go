package expenseus_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/saifahn/expenseus"
	"github.com/saifahn/expenseus/redis"
	"github.com/stretchr/testify/assert"
)

func TestCreatingExpensesAndRetrievingThem(t *testing.T) {
	wantedExpenseDetails := []expenseus.ExpenseDetails{
		expenseus.TestSeanExpenseDetails,
		expenseus.TestTomomiExpenseDetails,
		expenseus.TestTomomiExpense2Details,
	}

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("error starting up miniredis, %v", err)
	}

	db := redis.New(mr.Addr())
	oauth := &expenseus.StubOauthConfig{}
	auth := &expenseus.StubSessionManager{}
	images := &expenseus.StubImageStore{}
	webservice := expenseus.NewWebService(db, oauth, auth, "", images)
	router := expenseus.InitRouter(webservice)

	// CREATE users in the db
	testUsers := []expenseus.User{expenseus.TestSeanUser, expenseus.TestTomomiUser}
	for _, u := range testUsers {
		userJSON, err := json.Marshal(u)
		if err != nil {
			t.Fatalf(err.Error())
		}

		response := httptest.NewRecorder()
		request := expenseus.NewCreateUserRequest(userJSON)
		request.AddCookie(&expenseus.ValidCookie)
		router.ServeHTTP(response, request)
	}

	// GET all users
	response := httptest.NewRecorder()
	request := expenseus.NewGetAllUsersRequest()
	request.AddCookie(&expenseus.ValidCookie)
	if err != nil {
		t.Fatalf("request could not be created, %v", err)
	}
	router.ServeHTTP(response, request)

	var usersGot []expenseus.User
	err = json.NewDecoder(response.Body).Decode(&usersGot)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	}

	// ASSERT that the users received are correct
	assert.Len(t, usersGot, len(testUsers))
	assert.ElementsMatch(t, usersGot, testUsers)

	// CREATE expenses in the db
	for _, ed := range wantedExpenseDetails {
		values := map[string]io.Reader{
			"userID":      strings.NewReader(ed.UserID),
			"expenseName": strings.NewReader(ed.Name),
		}
		request := expenseus.NewCreateExpenseRequest(values)
		request.AddCookie(&expenseus.ValidCookie)
		router.ServeHTTP(httptest.NewRecorder(), request)
	}

	// GET all expenses
	response = httptest.NewRecorder()
	request = expenseus.NewGetAllExpensesRequest()
	request.AddCookie(&expenseus.ValidCookie)
	router.ServeHTTP(response, request)

	var expensesGot []expenseus.Expense
	err = json.NewDecoder(response.Body).Decode(&expensesGot)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	}

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, len(wantedExpenseDetails), len(expensesGot))
	// ASSERT there is an expense with the right details
	for _, expense := range expensesGot {
		assert.Contains(t, wantedExpenseDetails, expense.ExpenseDetails)
	}

	// GET one user's expenses
	response = httptest.NewRecorder()
	request = expenseus.NewGetExpensesByUsernameRequest(expenseus.TestTomomiUser.Username)
	request.AddCookie(&expenseus.ValidCookie)
	router.ServeHTTP(response, request)

	// reset expensesGot to an empty slice
	expensesGot = []expenseus.Expense{}
	err = json.NewDecoder(response.Body).Decode(&expensesGot)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	}

	assert.Len(t, expensesGot, 2)
	// ASSERT the expense details are the same. The ID is assigned separately, and we're currently not going to test that implementation
	var edsGot []expenseus.ExpenseDetails
	for _, e := range expensesGot {
		edsGot = append(edsGot, e.ExpenseDetails)
	}
	assert.Contains(t, edsGot, expenseus.TestTomomiExpense.ExpenseDetails)
	assert.Contains(t, edsGot, expenseus.TestTomomiExpense2.ExpenseDetails)
}

func TestRestrictedRoutesAndGettingSelf(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("error starting up miniredis")
	}

	db := redis.New(mr.Addr())
	oauth := &expenseus.StubOauthConfig{}
	auth := &expenseus.StubSessionManager{}
	images := &expenseus.StubImageStore{}
	webservice := expenseus.NewWebService(db, oauth, auth, "", images)
	router := expenseus.InitRouter(webservice)

	// try to create a user
	testUser := expenseus.User{
		Username: "testUser",
		Name:     "Test User",
		ID:       "test_user",
	}
	testUserJSON, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf(err.Error())
	}

	response := httptest.NewRecorder()
	request := expenseus.NewCreateUserRequest(testUserJSON)
	router.ServeHTTP(response, request)
	// with no valid cookie, should be unauthorized
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	// try again with a valid cookie
	response = httptest.NewRecorder()
	request = expenseus.NewCreateUserRequest(testUserJSON)
	request.AddCookie(&expenseus.ValidCookie)
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusAccepted, response.Code)

	// try and get a user from a cookie that doesn't exist in the db
	response = httptest.NewRecorder()
	request = expenseus.NewGetSelfRequest()
	request.AddCookie(&expenseus.ValidCookie)
	router.ServeHTTP(response, request)

	// sessionCookie, _ := request.Cookie("session")
	// println(sessionCookie.Value)

	assert.Equal(t, http.StatusNotFound, response.Code)

	// try and get the user created earlier
	response = httptest.NewRecorder()
	request = expenseus.NewGetSelfRequest()
	request.AddCookie(&http.Cookie{
		Name:  expenseus.ValidCookie.Name,
		Value: testUser.ID,
	})
	router.ServeHTTP(response, request)

	var got expenseus.User

	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, testUser, got)
}
