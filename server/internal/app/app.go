package app

import (
	"context"
	"mime/multipart"
	"net/http"

	"golang.org/x/oauth2"
)

type contextKey int

const (
	jsonContentType  = "application/json"
	SessionCookieKey = "expenseus-session"
)

type Store interface {
	GetExpense(id string) (Expense, error)
	GetExpensesByUsername(username string) ([]Expense, error)
	GetAllExpenses() ([]Expense, error)
	CreateExpense(expenseDetails ExpenseDetails) error
	CreateUser(user User) error
	GetUser(id string) (User, error)
	GetAllUsers() ([]User, error)
}

type Auth interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	GetInfoAndGenerateUser(state string, code string) (User, error)
}

type SessionManager interface {
	Validate(r *http.Request) bool
	Save(rw http.ResponseWriter, r *http.Request)
	GetUserID(r *http.Request) (string, error)
	Remove(rw http.ResponseWriter, r *http.Request)
}

type ImageStore interface {
	Upload(file multipart.File, header multipart.FileHeader) (string, error)
	Validate(file multipart.File) (bool, error)
	AddImageToExpense(expense Expense) (Expense, error)
}

type App struct {
	store    Store
	auth     Auth
	sessions SessionManager
	images   ImageStore
	frontend string
}

func New(store Store, oauth Auth, sessions SessionManager, frontend string, images ImageStore) *App {
	return &App{
		store:    store,
		auth:     oauth,
		sessions: sessions,
		frontend: frontend,
		images:   images,
	}
}
