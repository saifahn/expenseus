package expenseus_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifahn/expenseus"
	"github.com/stretchr/testify/assert"
)

type InMemoryExpenseStore struct {
	expenses map[string]expenseus.Expense
}

func (s *InMemoryExpenseStore) GetExpense(id string) (expenseus.Expense, error) {
	expense := s.expenses[id]
	return expense, nil
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

func (s *InMemoryExpenseStore) GetAllExpenses() []expenseus.Expense {
	var expenses []expenseus.Expense
	for _, e := range s.expenses {
		expenses = append(expenses, e)
	}
	return expenses
}

func TestCreatingExpensesAndRetrievingThem(t *testing.T) {
	store := InMemoryExpenseStore{
		map[string]expenseus.Expense{},
	}
	webservice := expenseus.NewWebService(&store)

	router := expenseus.InitRouter(webservice)

	router.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest("tomomi", "test expense 01"))
	router.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest("sean", "test expense 02"))
	router.ServeHTTP(httptest.NewRecorder(), expenseus.NewCreateExpenseRequest("tomomi", "test expense 03"))

	response := httptest.NewRecorder()
	request := expenseus.NewGetExpensesByUserRequest("tomomi")
	router.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Body.String(), `[test expense 01 test expense 03]`)
}
