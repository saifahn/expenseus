package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/saifahn/expenseus/internal/app"
)

func Init(a *app.App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
		r.Group(func(r chi.Router) {
			r.Use(a.VerifyUser)

			r.Post("/transactions", a.CreateTransaction)
			r.With(userIDCtx).
				Get("/transactions/user/{userID}", a.GetTransactionsByUser)
			r.With(userIDCtx).
				With(dateRangeCtx).
				Get("/transactions/user/{userID}/all", a.GetAllTxnsByUserBetweenDates)
			r.With(userIDCtx).
				With(dateRangeCtx).
				Get("/transactions/user/{userID}/range", a.GetTxnsBetweenDates)
			r.With(transactionIDCtx).
				Get("/transactions/{transactionID}", a.GetTransaction)
			r.With(transactionIDCtx).
				Put("/transactions/{transactionID}", a.UpdateTransaction)
			r.With(transactionIDCtx).
				Delete("/transactions/{transactionID}", a.DeleteTransaction)
			r.Post("/transactions/shared/settle", a.SettleTxns)

			r.Get("/users", a.ListUsers)
			r.Get("/users/self", a.GetSelf)
			r.Post("/users", a.CreateUser)
			r.With(userIDCtx).
				Get("/users/{userID}", a.GetUser)

			r.Post("/trackers", a.CreateTracker)
			r.With(trackerIDCtx).
				Get("/trackers/{trackerID}", a.GetTrackerByID)
			r.With(trackerIDCtx).
				Get("/trackers/{trackerID}/transactions", a.GetTxnsByTracker)
			r.With(trackerIDCtx).
				Post("/trackers/{trackerID}/transactions", a.CreateSharedTxn)
			r.With(trackerIDCtx).
				With(dateRangeCtx).
				Get("/trackers/{trackerID}/transactions/range", a.GetTxnsByTrackerBetweenDates)
			r.With(transactionIDCtx).
				With(trackerIDCtx).
				Put("/trackers/{trackerID}/transactions/{transactionID}", a.UpdateSharedTxn)
			r.With(trackerIDCtx).
				Delete("/trackers/{trackerID}/transactions/{transactionID}", a.DeleteSharedTxn)
			r.With(trackerIDCtx).
				Get("/trackers/{trackerID}/transactions/unsettled", a.GetUnsettledTxnsByTracker)
			r.With(userIDCtx).
				Get("/trackers/user/{userID}", a.GetTrackersByUser)
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

// Gets the TrackerID from the URL and adds it to the TrackerID context for the request.
func trackerIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "trackerID")
		ctx := context.WithValue(r.Context(), app.CtxKeyTrackerID, id)
		ctx = context.WithValue(ctx, app.CtxKeyUserID, r.Context().Value(app.CtxKeyUserID))
		ctx = context.WithValue(ctx, app.CtxKeyTransactionID, r.Context().Value(app.CtxKeyTransactionID))
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func dateRangeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		if from == "" || to == "" {
			http.Error(w, "from and to must both be supplied", http.StatusBadRequest)
		}
		fromInt64, err := strconv.ParseInt(from, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("from was not a valid number: %v", err.Error()), http.StatusBadRequest)
		}
		toInt64, err := strconv.ParseInt(to, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("to was not a valid number: %v", err.Error()), http.StatusBadRequest)
		}
		ctx := context.WithValue(r.Context(), app.CtxKeyDateFrom, fromInt64)
		ctx = context.WithValue(ctx, app.CtxKeyDateTo, toInt64)
		ctx = context.WithValue(ctx, app.CtxKeyUserID, r.Context().Value(app.CtxKeyUserID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
