package expenseus_test

import (
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

	// create expenses in the db
	for _, ed := range wantedExpenseDetails {
		router.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest(
			ed.UserID, ed.Name,
		))
	}

	// get all expenses
	response := httptest.NewRecorder()
	request := expenseus.NewGetAllExpensesRequest()
	router.ServeHTTP(response, request)

	var got []expenseus.Expense
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	}

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, len(wantedExpenseDetails), len(got))
	// assert there is an expense with the right details
	for _, expense := range got {
		assert.Contains(t, wantedExpenseDetails, expense.ExpenseDetails)
	}

	// get one user's expenses
	// response = httptest.NewRecorder()
	// request = expenseus.NewGetExpensesByUsernameRequest(expenseus.TestTomomiUser.Username)
	// router.ServeHTTP(response, request)

	// // reset got to an empty slice
	// got = []expenseus.Expense{}
	// err = json.NewDecoder(response.Body).Decode(&got)
	// if err != nil {
	// 	t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	// }

	// assert.Len(t, got, 2)
	// assert.Contains(t, got, expenseus.TestTomomiExpense)
	// assert.Contains(t, got, expenseus.TestTomomiExpense2)
}
