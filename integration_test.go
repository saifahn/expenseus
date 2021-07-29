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
	wantedExpenses := []expenseus.Expense{
		{"tomomi", "test expense 01"},
		{"sean", "test expense 02"},
		{"tomomi", "test expense 03"},
	}

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("error starting up miniredis, %v", err)
	}

	db := redis.New(mr.Addr())
	webservice := expenseus.NewWebService(db)

	router := expenseus.InitRouter(webservice)

	for _, e := range wantedExpenses {
		router.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest(
			e.User, e.Name,
		))
	}

	response := httptest.NewRecorder()
	request := expenseus.NewGetAllExpensesRequest()

	router.ServeHTTP(response, request)

	var got []expenseus.Expense
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	}
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, len(wantedExpenses), len(got))
	assert.ElementsMatch(t, wantedExpenses, got)

	// assert.Contains(t, got, expenseus.Expense{User: "tomomi", Name: "test expense 01"})
	// assert.Contains(t, got, expenseus.Expense{User: "tomomi", Name: "test expense 03"})
}
