package expenseus

import (
	"encoding/json"
	"net/http"
)

type contextKey int

const (
	CtxKeyExpenseID contextKey = iota
	CtxKeyUsername  contextKey = iota
)

type ExpenseStore interface {
	GetExpense(id string) (Expense, error)
	GetExpensesByUser(user string) ([]Expense, error)
	GetAllExpenses() ([]Expense, error)
	RecordExpense(expense Expense) error
}

type Expense struct {
	Name string `json:"name"`
	User string `json:"user"`
}

func NewWebService(store ExpenseStore) *WebService {
	return &WebService{store}
}

type WebService struct {
	store ExpenseStore
}

// GetExpense handles a HTTP request to get an expense by ID, returning the expense.
func (wb *WebService) GetExpense(rw http.ResponseWriter, r *http.Request) {
	expenseID := r.Context().Value(CtxKeyExpenseID).(string)

	expense, err := wb.store.GetExpense(expenseID)

	// TODO: should account for different kinds of errors
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	rw.Header().Set("content-type", "application/json")
	json.NewEncoder(rw).Encode(expense)
}

// GetExpensesByUser handles a HTTP request to get all expenses of a user,
// returning a list of expenses.
func (wb *WebService) GetExpensesByUser(rw http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(CtxKeyUsername).(string)

	expenses, err := wb.store.GetExpensesByUser(username)

	// TODO: account for different errors
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: handle error
	json.NewEncoder(rw).Encode(expenses)

	rw.WriteHeader(http.StatusOK)
}

// GetAllExpenses handles a HTTP request to get all expenses, returning a list
// of expenses.
func (wb *WebService) GetAllExpenses(rw http.ResponseWriter, r *http.Request) {
	expenses, err := wb.store.GetAllExpenses()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: handle error
	json.NewEncoder(rw).Encode(expenses)

	rw.WriteHeader(http.StatusOK)
}

// CreateExpense handles a HTTP request to create a new expense.
func (wb *WebService) CreateExpense(rw http.ResponseWriter, r *http.Request) {
	var e Expense
	err := json.NewDecoder(r.Body).Decode(&e)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = wb.store.RecordExpense(e)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusAccepted)
}
