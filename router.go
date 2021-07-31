package expenseus

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitRouter(wb *WebService) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/expenses", func(r chi.Router) {
		r.Get("/", wb.GetAllExpenses)
		r.Post("/", wb.CreateExpense)

		r.Route("/user/{username}", func(r chi.Router) {
			r.Use(UsernameCtx)
			r.Get("/", wb.GetExpensesByUsername)
		})

		r.Route("/{expenseID}", func(r chi.Router) {
			r.Use(ExpenseIDCtx)
			r.Get("/", wb.GetExpense)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", wb.ListUsers)
		r.Post("/", wb.CreateUser)
	})

	return r
}

// Gets the ID from the URL and adds it to the id context for the request.
func ExpenseIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		expenseID := chi.URLParam(r, "expenseID")
		ctx := context.WithValue(r.Context(), CtxKeyExpenseID, expenseID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// Gets the username from the URL and adds it to the user context for the request.
func UsernameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		ctx := context.WithValue(r.Context(), CtxKeyUsername, username)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
