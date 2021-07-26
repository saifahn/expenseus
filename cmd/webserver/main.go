package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()

	r.Route("/expenses", func(r chi.Router) {
		r.Get("/", wb.GetAllExpenses)
		r.Post("/", wb.CreateExpense)

		r.Route("/users/{username}", func(r chi.Router) {
			r.Use(UsernameCtx)
			r.Get("/", wb.GetExpensesByUser)
		})

		r.Route("/{expenseID}", func(r chi.Router) {
			r.Use(ExpenseIDCtx)
			r.Get("/", wb.GetExpense)
		})
	})

	log.Fatal(http.ListenAndServe(":5000", r))
}

// Gets the ID from the URL and adds it to the id context for the request.
func ExpenseIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		expenseID := chi.URLParam(r, "expenseID")
		ctx := context.WithValue(r.Context(), expenseus.CtxKeyExpenseID, expenseID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// Gets the username from the URL and adds it to the user context for the request.
func UsernameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		ctx := context.WithValue(r.Context(), expenseus.CtxKeyUsername, username)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
