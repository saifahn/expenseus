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

func (s *InMemoryExpenseStore) GetExpenseNameByID(id string) string {
	return "not implemented yet"
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

func (s *InMemoryExpenseStore) RecordExpense(expense expenseus.Expense) {
	s.Expenses = append(s.Expenses, expense)
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
		r.Post("/", wb.CreateExpense)

		r.Route("/users/{username}", func(r chi.Router) {
			r.Use(UsernameCtx)
			r.Get("/", wb.GetExpensesByUser)
		})

		r.Route("/{expenseID}", func(r chi.Router) {
			r.Use(ExpenseIDCtx)
			r.Get("/", wb.GetExpenseByID)
		})
	})

	log.Fatal(http.ListenAndServe(":5000", r))
}

// Gets the ID from the URL and adds it to the id context for the request.
func ExpenseIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		expenseID := chi.URLParam(r, "expenseID")
		ctx := context.WithValue(r.Context(), "expenseID", expenseID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// Gets the username from the URL and adds it to the user context for the request.
func UsernameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		ctx := context.WithValue(r.Context(), "username", username)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
