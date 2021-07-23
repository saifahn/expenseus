package expenseus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExpenseStore interface {
	GetExpenseNameById(id string) string
	GetExpenseNamesByUser(user string) []string
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

func (wb *WebService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()
	router.HandleFunc("/expenses/user/", wb.GetExpensesByUser)
	router.HandleFunc("/expenses/", wb.GetExpenseByID)
	router.HandleFunc("/expenses", wb.createExpenseHandler)

	router.ServeHTTP(w, r)
}

// GetExpenseByID handles a HTTP request to get an expense by ID, returning the expense name.
// TODO: update the comment when you return the expense completely
func (wb *WebService) GetExpenseByID(rw http.ResponseWriter, r *http.Request) {
	expenseId := r.Context().Value("id").(string)

	expenseName := wb.store.GetExpenseNameById(expenseId)

	if expenseName == "" {
		rw.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(rw, expenseName)
}

// GetExpensesByUser handles a HTTP request to get all expenses of a user,
// returning a list of expense names.
// TODO: update this comment
func (wb *WebService) GetExpensesByUser(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(string)

	expenses := wb.store.GetExpenseNamesByUser(user)

	fmt.Fprint(rw, expenses)
}

func (wb *WebService) createExpenseHandler(rw http.ResponseWriter, r *http.Request) {
	var e Expense
	err := json.NewDecoder(r.Body).Decode(&e)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	wb.store.RecordExpense(e)
	rw.WriteHeader(http.StatusAccepted)
}
