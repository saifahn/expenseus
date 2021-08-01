package expenseus_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	webservice := expenseus.NewWebService(db)
	router := expenseus.InitRouter(webservice)

	// CREAT users in the db
	testUsers := []expenseus.User{expenseus.TestSeanUser, expenseus.TestTomomiUser}
	for _, u := range testUsers {
		userJSON, err := json.Marshal(u)
		if err != nil {
			t.Fatalf(err.Error())
		}

		response := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatalf("request could not be created, %v", err)
		}
		router.ServeHTTP(response, request)
	}

	// GET all users
	response := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/users", nil)
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
		router.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest(
			ed.UserID, ed.Name,
		))
	}

	// GET all expenses
	response = httptest.NewRecorder()
	request = expenseus.NewGetAllExpensesRequest()
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
