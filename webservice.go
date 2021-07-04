package expenseus

import (
	"fmt"
	"net/http"
	"strings"
)

func WebService(w http.ResponseWriter, r *http.Request) {
	expenseId := strings.TrimPrefix(r.URL.Path, "/expenses/")

	if expenseId == "9281" {
		fmt.Fprint(w, "Expense 9281")
		return
	}

	if expenseId == "1" {
		fmt.Fprint(w, "Expense 1")
		return
	}

}
