package expenseus

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

type contextKey int

const (
	CtxKeyExpenseID contextKey = iota
	CtxKeyUsername  contextKey = iota
	jsonContentType            = "application/json"
)

type ExpenseStore interface {
	GetExpense(id string) (Expense, error)
	GetExpensesByUsername(username string) ([]Expense, error)
	GetAllExpenses() ([]Expense, error)
	RecordExpense(expenseDetails ExpenseDetails) error
	CreateUser(user User) error
	GetAllUsers() ([]User, error)
}

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

type ExpenseDetails struct {
	Name   string `json:"name"`
	UserID string `json:"userid"`
}

type Expense struct {
	ExpenseDetails
	ID string `json:"id"`
}

type GoogleUserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
}

type GoogleOauthConfig interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	getUserInfo(state string, code string) (GoogleUserInfo, error)
}

type WebService struct {
	store       ExpenseStore
	oauthConfig GoogleOauthConfig
}

func NewWebService(store ExpenseStore, oauth GoogleOauthConfig) *WebService {
	return &WebService{store: store, oauthConfig: oauth}
}

func (wb *WebService) OauthCallback(rw http.ResponseWriter, r *http.Request) {
	userinfo, err := wb.oauthConfig.getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a user with default information from the oauth user info
	user := User{
		Username: userinfo.Email,
		Name:     userinfo.Name,
		ID:       userinfo.ID,
	}

	wb.store.CreateUser(user)
}

// GetExpense handles a HTTP request to get an expense by ID, returning the expense.
func (wb *WebService) GetExpense(rw http.ResponseWriter, r *http.Request) {
	expenseID := r.Context().Value(CtxKeyExpenseID).(string)

	expense, err := wb.store.GetExpense(expenseID)

	// TODO: should account for different kinds of errors
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(expense)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetExpensesByUsername handles a HTTP request to get all expenses of a user,
// returning a list of expenses.
func (wb *WebService) GetExpensesByUsername(rw http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(CtxKeyUsername).(string)

	expenses, err := wb.store.GetExpensesByUsername(username)

	// TODO: account for different errors
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(expenses)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAllExpenses handles a HTTP request to get all expenses, returning a list
// of expenses.
func (wb *WebService) GetAllExpenses(rw http.ResponseWriter, r *http.Request) {
	expenses, err := wb.store.GetAllExpenses()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(expenses)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateExpense handles a HTTP request to create a new expense.
func (wb *WebService) CreateExpense(rw http.ResponseWriter, r *http.Request) {
	var ed ExpenseDetails
	err := json.NewDecoder(r.Body).Decode(&ed)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = wb.store.RecordExpense(ed)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

func (wb *WebService) CreateUser(rw http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = wb.store.CreateUser(u)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

func (wb *WebService) ListUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := wb.store.GetAllUsers()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
