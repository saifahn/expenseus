package expenseus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type contextKey int

const (
	CtxKeyExpenseID contextKey = iota
	CtxKeyUsername  contextKey = iota
)

type ExpenseStore interface {
	GetExpense(id string) (Expense, error)
	GetExpenseNamesByUser(user string) []string
	GetAllExpenses() []Expense
	RecordExpense(expense Expense)
}

type Expense struct {
	Name string
	User string
}

func NewWebService(store ExpenseStore) *WebService {
	return &WebService{store}
}

type WebService struct {
	store ExpenseStore
}

// GetExpense handles a HTTP request to get an expense by ID, returning the expense.
func (wb *WebService) GetExpense(rw http.ResponseWriter, r *http.Request) {
	expenseId := r.Context().Value(CtxKeyExpenseID).(string)

	expense, err := wb.store.GetExpense(expenseId)

	// TODO: should account for different kinds of errors
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	rw.Header().Set("content-type", "application/json")
	json.NewEncoder(rw).Encode(expense)
}

// GetExpensesByUser handles a HTTP request to get all expenses of a user,
// returning a list of expense names.
// TODO: update this comment
func (wb *WebService) GetExpensesByUser(rw http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(CtxKeyUsername).(string)

	expenses := wb.store.GetExpenseNamesByUser(username)

	fmt.Fprint(rw, expenses)
}

// GetAllExpenses handles a HTTP request to get all expenses, return a list of
// expense names.
func (wb *WebService) GetAllExpenses(rw http.ResponseWriter, r *http.Request) {
	expenses := wb.store.GetAllExpenses()

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

	wb.store.RecordExpense(e)
	rw.WriteHeader(http.StatusAccepted)
}
