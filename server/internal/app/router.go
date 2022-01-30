package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRouter(wb *WebService) *chi.Mux {
	r := chi.NewRouter()

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		// TODO: use environment variables to determine allowed origins
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	fs := http.FileServer(http.Dir("./web/dist"))
	r.Handle("/*", fs)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/expenses", func(r chi.Router) {
			r.Use(wb.VerifyUser)
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
			r.Use(wb.VerifyUser)
			r.Get("/", wb.ListUsers)
			r.Post("/", wb.CreateUser)
			r.Get("/self", wb.GetSelf)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(UserIDCtx)
				r.Get("/", wb.GetUser)
			})
		})

		r.Get("/login_google", wb.OauthLogin)
		r.Get("/callback_google", wb.OauthCallback)
		r.Get("/logout", wb.LogOut)
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

// Gets the UserID from the URL and adds it to the UserID context for the request.
func UserIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "userID")
		ctx := context.WithValue(r.Context(), CtxKeyUserID, id)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
