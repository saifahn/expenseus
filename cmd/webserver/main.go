package main

import (
	"log"
	"net/http"

	"github.com/saifahn/expenseus"
)

type InMemoryExpenseStore struct {
	Expenses []expenseus.Expense
}

func (s *InMemoryExpenseStore) GetExpense(id string) (expenseus.Expense, error) {
	expense := &expenseus.Expense{}
	return *expense, nil
}

func (s *InMemoryExpenseStore) GetExpenseNamesByUser(user string) []string {
	var expenseNames []string
	for _, e := range s.Expenses {
		if e.User == user {
			expenseNames = append(expenseNames, e.Name)
		}
	}
	return expenseNames
}

func (s *InMemoryExpenseStore) GetAllExpenseNames() []string {
	var expenseNames []string
	for _, e := range s.Expenses {
		expenseNames = append(expenseNames, e.Name)
	}
	return expenseNames
}

func (s *InMemoryExpenseStore) RecordExpense(expense expenseus.Expense) {
	s.Expenses = append(s.Expenses, expense)
}

func (s *InMemoryExpenseStore) GetAllExpenses() []expenseus.Expense {
	return s.Expenses
}

func main() {
	wb := expenseus.NewWebService(&InMemoryExpenseStore{
		Expenses: []expenseus.Expense{
			{
				Name: "tomomi-01",
				User: "tomomi",
			},
			{
				Name: "tomomi-02",
				User: "tomomi",
			},
			{
				Name: "tomomi-03",
				User: "tomomi",
			},
			{
				Name: "sean-01",
				User: "sean",
			},
		},
	})

	r := expenseus.InitRouter(wb)

	log.Fatal(http.ListenAndServe(":5000", r))
}
