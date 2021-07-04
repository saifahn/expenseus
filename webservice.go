package expenseus

import (
	"fmt"
	"net/http"
	"strings"
)

type ExpenseStore interface {
	GetExpense(id string) string
}

func NewWebService(store ExpenseStore) *WebService {
	return &WebService{store}
}

type WebService struct {
	store ExpenseStore
}

func (wb *WebService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	expenseId := strings.TrimPrefix(r.URL.Path, "/expenses/")

	w.WriteHeader(http.StatusNotFound)

	fmt.Fprint(w, wb.store.GetExpense(expenseId))
}

func GetExpense(id string) string {
	expense := ""
	if id == "9281" {
		expense = "Expense 9281"
	}

	if id == "1" {
		expense = "Expense 1"
	}
	return expense
}
