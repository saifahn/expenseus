package expenseus_test

import (
	"encoding/json"
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

func (s *InMemoryExpenseStore) GetExpensesByUser(user string) ([]expenseus.Expense, error) {
	var expenses []expenseus.Expense
	for _, e := range s.expenses {
		if e.User == user {
			expenses = append(expenses, e)
		}
	}
	return expenses, nil
}

func (s *InMemoryExpenseStore) RecordExpense(e expenseus.Expense) error {
	testId := fmt.Sprintf("tid-%v", e.Name)
	s.expenses[testId] = e
	return nil
}

func (s *InMemoryExpenseStore) GetAllExpenses() ([]expenseus.Expense, error) {
	var expenses []expenseus.Expense
	for _, e := range s.expenses {
		expenses = append(expenses, e)
	}
	return expenses, nil
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

	var got []expenseus.Expense
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	}

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, 2, len(got))
	assert.Contains(t, got, expenseus.Expense{User: "tomomi", Name: "test expense 01"})
	assert.Contains(t, got, expenseus.Expense{User: "tomomi", Name: "test expense 03"})
}
