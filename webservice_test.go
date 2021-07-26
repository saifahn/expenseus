package expenseus

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSeanExpense = Expense{
	Name: "Expense 1",
	User: "sean",
}

var testTomomiExpense = Expense{
	Name: "Expense 9281",
	User: "tomomi",
}

func TestGetExpenseByID(t *testing.T) {
	store := StubExpenseStore{
		map[string]Expense{
			"1":    testSeanExpense,
			"9281": testTomomiExpense,
		},
	}
	webservice := &WebService{&store}

	t.Run("get an expense by id", func(t *testing.T) {
		request := NewGetExpenseRequest("1")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, response.Result().Header.Get("content-type"), "application/json")
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, got, testSeanExpense)
	})

	t.Run("gets another expense by id", func(t *testing.T) {
		request := NewGetExpenseRequest("9281")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, response.Result().Header.Get("content-type"), "application/json")
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, got, testTomomiExpense)
	})

	t.Run("returns 404 on non-existent expense", func(t *testing.T) {
		request := NewGetExpenseRequest("13371337")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestGetExpenseByUser(t *testing.T) {
	store := StubExpenseStore{
		map[string]Expense{
			"1":    testSeanExpense,
			"9281": testTomomiExpense,
		},
	}
	webservice := NewWebService(&store)

	t.Run("gets Tomomi's expenses", func(t *testing.T) {
		request := NewGetExpensesByUserRequest("tomomi")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpensesByUser)
		handler.ServeHTTP(response, request)

		AssertResponseStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), fmt.Sprintf("[%v]", testTomomiExpense.Name))
	})

	t.Run("gets Sean's expenses", func(t *testing.T) {
		request := NewGetExpensesByUserRequest("sean")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpensesByUser)
		handler.ServeHTTP(response, request)

		AssertResponseStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), fmt.Sprintf("[%v]", testSeanExpense.Name))
	})
}

func TestCreateExpense(t *testing.T) {
	store := StubExpenseStore{
		map[string]Expense{},
	}
	webservice := NewWebService(&store)

	t.Run("creates a new expense on POST", func(t *testing.T) {
		request := NewCreateExpenseRequest("tomomi", "Test Expense")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.CreateExpense)
		handler.ServeHTTP(response, request)

		AssertResponseStatus(t, response.Code, http.StatusAccepted)

		if len(store.expenses) != 1 {
			t.Errorf("got %d expenses, want %d", len(store.expenses), 1)
		}
	})
}

func TestGetAllExpenses(t *testing.T) {
	t.Run("gets all expenses with one expense", func(t *testing.T) {
		wantedExpenses := []Expense{
			{User: "tomomi", Name: "test expense 01"},
		}
		store := StubExpenseStore{
			map[string]Expense{
				"01": {
					User: "tomomi",
					Name: "test expense 01",
				},
			},
		}
		webservice := NewWebService(&store)

		req, err := http.NewRequest(http.MethodGet, "/expenses", nil)
		if err != nil {
			t.Errorf("there was an error in creating the request")
		}
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetAllExpenses)
		handler.ServeHTTP(response, req)

		var got []Expense
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		AssertResponseStatus(t, response.Code, http.StatusOK)
		if !reflect.DeepEqual(got, wantedExpenses) {
			t.Errorf("got expenses %v, wanted %v", got, wantedExpenses)
		}
	})

	t.Run("gets all expenses with more than one expense", func(t *testing.T) {
		wantedExpenses := []Expense{
			{User: "tomomi", Name: "test expense 01"},
			{User: "sean", Name: "test expense 02"},
			{User: "tomomi", Name: "test expense 03"},
		}
		store := StubExpenseStore{
			map[string]Expense{
				"01": {
					User: "tomomi",
					Name: "test expense 01",
				},
				"02": {
					User: "sean",
					Name: "test expense 02",
				},
				"03": {
					User: "tomomi",
					Name: "test expense 03",
				},
			},
		}
		webservice := NewWebService(&store)
		req, err := http.NewRequest(http.MethodGet, "/expenses", nil)
		if err != nil {
			t.Errorf("there was an error in creating the request")
		}
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetAllExpenses)
		handler.ServeHTTP(response, req)

		var got []Expense
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		AssertResponseStatus(t, response.Code, http.StatusOK)
		assert.Equal(t, len(got), len(wantedExpenses))
		assert.ElementsMatch(t, got, wantedExpenses)
	})

}

type StubExpenseStore struct {
	expenses map[string]Expense
}

func (s *StubExpenseStore) GetExpense(id string) (Expense, error) {
	expense := s.expenses[id]
	// check for empty Expense
	if expense == (Expense{}) {
		return Expense{}, errors.New("expense not found")
	}
	return expense, nil
}

func (s *StubExpenseStore) GetExpenseNamesByUser(user string) []string {
	var expenseNames []string
	for _, e := range s.expenses {
		if e.User == user {
			expenseNames = append(expenseNames, e.Name)
		}
	}
	return expenseNames
}

func (s *StubExpenseStore) RecordExpense(e Expense) {
	testId := fmt.Sprintf("tid-%v", e.Name)
	s.expenses[testId] = e
}

func (s *StubExpenseStore) GetAllExpenses() []Expense {
	var expenses []Expense
	for _, e := range s.expenses {
		expenses = append(expenses, e)
	}
	return expenses
}
