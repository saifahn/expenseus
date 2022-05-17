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
	GetTransaction(txnID string) (Transaction, error)
	GetTransactionsByUser(userID string) ([]Transaction, error)
	GetAllTransactions() ([]Transaction, error)
	CreateTransaction(transactionDetails TransactionDetails) error
	CreateUser(user User) error
	GetUser(id string) (User, error)
	GetAllUsers() ([]User, error)
	CreateTracker(tracker Tracker) error
	GetTracker(trackerID string) (Tracker, error)
	GetTrackersByUser(userID string) ([]Tracker, error)
	GetTxnsByTracker(trackerID string) ([]SharedTransaction, error)
	CreateSharedTxn(txn SharedTransaction) error
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
