package expenseus

import (
	"fmt"
	"net/http"
	"strings"
)

func WebService(w http.ResponseWriter, r *http.Request) {
	expenseId := strings.TrimPrefix(r.URL.Path, "/expenses/")
	expense := GetExpense(expenseId)
	fmt.Fprint(w, expense)
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
