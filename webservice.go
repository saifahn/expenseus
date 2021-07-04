package expenseus

import (
	"fmt"
	"net/http"
	"strings"
)

type ExpenseStore interface {
	GetExpenseNameById(id string) string
}

type Expense struct {
	name string
	user string
}

func NewWebService(store ExpenseStore) *WebService {
	return &WebService{store}
}

type WebService struct {
	store ExpenseStore
}

func (wb *WebService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	expenseId := strings.TrimPrefix(r.URL.Path, "/expenses/")

	expenseName := wb.store.GetExpenseNameById(expenseId)

	if expenseName == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, expenseName)
}
