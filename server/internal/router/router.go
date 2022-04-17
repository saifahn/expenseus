package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/saifahn/expenseus/internal/app"
)

func Init(a *app.App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

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
		r.Route("/transactions", func(r chi.Router) {
			r.Use(a.VerifyUser)
			r.Get("/", a.GetAllTransactions)
			r.Post("/", a.CreateTransaction)

			r.Route("/user/{userID}", func(r chi.Router) {
				r.Use(userIDCtx)
				r.Get("/", a.GetTransactionsByUser)
			})

			r.Route("/{transactionID}", func(r chi.Router) {
				r.Use(transactionIDCtx)
				r.Get("/", a.GetTransaction)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(a.VerifyUser)
			r.Get("/", a.ListUsers)
			r.Post("/", a.CreateUser)
			r.Get("/self", a.GetSelf)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(userIDCtx)
				r.Get("/", a.GetUser)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(a.VerifyUser)
			r.Post("/trackers", a.CreateTracker)
		})

		r.Get("/login_google", a.OauthLogin)
		r.Get("/callback_google", a.OauthCallback)
		r.Get("/logout", a.LogOut)
	})

	return r
}

// Gets the ID from the URL and adds it to the id context for the request.
func transactionIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		transactionID := chi.URLParam(r, "transactionID")
		ctx := context.WithValue(r.Context(), app.CtxKeyTransactionID, transactionID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// Gets the UserID from the URL and adds it to the UserID context for the request.
func userIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "userID")
		ctx := context.WithValue(r.Context(), app.CtxKeyUserID, id)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
