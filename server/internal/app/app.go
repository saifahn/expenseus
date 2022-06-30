package app

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
)

type contextKey int

const (
	CtxKeyTrackerID contextKey = iota
	CtxKeyTransactionID
	CtxKeyUserID
	CtxKeyDateFrom
	CtxKeyDateTo
)

const (
	jsonContentType  = "application/json"
	SessionCookieKey = "expenseus-session"
)

type Store interface {
	// Transactions
	GetTransaction(userID, txnID string) (Transaction, error)
	GetTransactionsByUser(userID string) ([]Transaction, error)
	GetTxnsBetweenDates(userID string, from, to int64) ([]Transaction, error)
	CreateTransaction(txn Transaction) error
	UpdateTransaction(txn Transaction) error
	DeleteTransaction(txnID, userID string) error

	// Users
	CreateUser(user User) error
	GetUser(id string) (User, error)
	GetAllUsers() ([]User, error)

	// SharedTxns
	CreateTracker(tracker Tracker) error
	GetTracker(trackerID string) (Tracker, error)
	GetTrackersByUser(userID string) ([]Tracker, error)
	CreateSharedTxn(txn SharedTransaction) error
	UpdateSharedTxn(txn SharedTransaction) error
	DeleteSharedTxn(input DelSharedTxnInput) error
	GetTxnsByTracker(trackerID string) ([]SharedTransaction, error)
	GetUnsettledTxnsByTracker(trackerID string) ([]SharedTransaction, error)
	SettleTxns(txns []SharedTransaction) error
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
	AddImageToTransaction(transaction Transaction) (Transaction, error)
}

type App struct {
	store    Store
	auth     Auth
	sessions SessionManager
	images   ImageStore
	frontend string
	validate *validator.Validate
}

func New(store Store, oauth Auth, sessions SessionManager, frontend string, images ImageStore) *App {
	validate := validator.New()

	return &App{
		store:    store,
		auth:     oauth,
		sessions: sessions,
		frontend: frontend,
		images:   images,
		validate: validate,
	}
}
