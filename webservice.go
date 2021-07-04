package expenseus

import (
	"fmt"
	"net/http"
	"strings"
)

type ExpenseStore interface {
	GetExpenseNameById(id string) string
	GetExpenseNamesByUser(user string) []string
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
	router := http.NewServeMux()

	router.Handle("/expenses/user/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := strings.TrimPrefix(r.URL.Path, "/expenses/user/")
		expenses := wb.store.GetExpenseNamesByUser(user)
		fmt.Fprint(rw, expenses)
	}))

	router.Handle("/expenses/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		expenseId := strings.TrimPrefix(r.URL.Path, "/expenses/")

		expenseName := wb.store.GetExpenseNameById(expenseId)

		if expenseName == "" {
			rw.WriteHeader(http.StatusNotFound)
		}

		fmt.Fprint(rw, expenseName)
	}))

	router.ServeHTTP(w, r)
}
