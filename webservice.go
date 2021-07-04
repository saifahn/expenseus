package expenseus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	router.HandleFunc("/expenses/user/", wb.expenseByUserHandler)
	router.HandleFunc("/expenses/", wb.expenseByIdHandler)
	router.HandleFunc("/expenses", wb.createExpenseHandler)

	router.ServeHTTP(w, r)
}

func (wb *WebService) expenseByIdHandler(rw http.ResponseWriter, r *http.Request) {
	expenseId := strings.TrimPrefix(r.URL.Path, "/expenses/")

	expenseName := wb.store.GetExpenseNameById(expenseId)

	if expenseName == "" {
		rw.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(rw, expenseName)
}

func (wb *WebService) expenseByUserHandler(rw http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/expenses/user/")

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
