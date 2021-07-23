package expenseus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testSeanExpense = Expense{
	Name: "Expense 1",
	User: "sean",
}

var testTomomiExpense = Expense{
	Name: "Expense 9281",
	User: "tomomi",
}

func TestGetExpenseById(t *testing.T) {
	store := StubExpenseStore{
		map[string]Expense{
			"1":    testSeanExpense,
			"9281": testTomomiExpense,
		},
	}
	webservice := &WebService{&store}

	t.Run("get an expense by id", func(t *testing.T) {
		request := NewGetExpenseByIdRequest("1")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.expenseByIdHandler)
		handler.ServeHTTP(response, request)

		AssertResponseStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), testSeanExpense.Name)
	})

	t.Run("gets another expense by id", func(t *testing.T) {
		request := NewGetExpenseByIdRequest("9281")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.expenseByIdHandler)
		handler.ServeHTTP(response, request)

		AssertResponseStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), testTomomiExpense.Name)
	})

	t.Run("returns 404 on non-existent expense", func(t *testing.T) {
		request := NewGetExpenseByIdRequest("13371337")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.expenseByIdHandler)
		handler.ServeHTTP(response, request)

		AssertResponseStatus(t, response.Code, http.StatusNotFound)
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
		request, _ := http.NewRequest(http.MethodGet, "/expenses/user/tomomi", nil)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.expenseByUserHandler)
		handler.ServeHTTP(response, request)

		AssertResponseStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), fmt.Sprintf("[%v]", testTomomiExpense.Name))
	})

	t.Run("gets Sean's expenses", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/expenses/user/sean", nil)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.expenseByUserHandler)
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

		handler := http.HandlerFunc(webservice.createExpenseHandler)
		handler.ServeHTTP(response, request)

		AssertResponseStatus(t, response.Code, http.StatusAccepted)

		if len(store.expenses) != 1 {
			t.Errorf("got %d expenses, want %d", len(store.expenses), 1)
		}
	})
}

type StubExpenseStore struct {
	expenses map[string]Expense
}

func (s *StubExpenseStore) GetExpenseNameById(id string) string {
	expense := s.expenses[id]
	return expense.Name
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
