package expenseus_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifahn/expenseus"
)

type InMemoryExpenseStore struct {
	expenses map[string]expenseus.Expense
}

func (s *InMemoryExpenseStore) GetExpenseNameByID(id string) string {
	expense := s.expenses[id]
	return expense.Name
}

func (s *InMemoryExpenseStore) GetExpenseNamesByUser(user string) []string {
	var expenseNames []string
	for _, e := range s.expenses {
		if e.User == user {
			expenseNames = append(expenseNames, e.Name)
		}
	}
	return expenseNames
}

func (s *InMemoryExpenseStore) RecordExpense(e expenseus.Expense) {
	testId := fmt.Sprintf("tid-%v", e.Name)
	s.expenses[testId] = e
}

func (s *InMemoryExpenseStore) GetAllExpenseNames() []string {
	var expenseNames []string
	for _, e := range s.expenses {
		expenseNames = append(expenseNames, e.Name)
	}
	return expenseNames
}

func TestCreatingExpensesAndRetrievingThem(t *testing.T) {
	store := InMemoryExpenseStore{
		map[string]expenseus.Expense{},
	}
	webservice := expenseus.NewWebService(&store)

	webservice.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest("tomomi", "test expense 01"))
	webservice.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest("sean", "test expense 02"))
	webservice.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest("tomomi", "test expense 03"))

	response := httptest.NewRecorder()
	request := expenseus.NewGetExpensesByUserRequest("tomomi")
	webservice.ServeHTTP(response, request)
	expenseus.AssertResponseStatus(t, response.Code, http.StatusOK)

	expenseus.AssertResponseBody(t, response.Body.String(), `[test expense 01 test expense 03]`)
}
