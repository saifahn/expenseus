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

	expense := wb.store.GetExpense(expenseId)

	if expense == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, expense)
}
