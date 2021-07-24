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
	return []string{"test", "something"}
}

func (s *InMemoryExpenseStore) RecordExpense(expense expenseus.Expense) {
	s.Expenses = append(s.Expenses, expense)
}

func main() {
	wb := expenseus.NewWebService(&InMemoryExpenseStore{})

	r := chi.NewRouter()

	r.Route("/expenses", func(r chi.Router) {

		r.Route("/{expenseID}", func(r chi.Router) {
			r.Use(IDCtx)
			r.Get("/", wb.GetExpenseByID)
		})
	})

	log.Fatal(http.ListenAndServe(":5000", r))
}

// Gets the ID from the URL and adds it to the id context for the request.
func IDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		expenseID := chi.URLParam(r, "expenseID")
		ctx := context.WithValue(r.Context(), "expenseID", expenseID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
