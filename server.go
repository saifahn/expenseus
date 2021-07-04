package expenseus

import (
	"fmt"
	"net/http"
)

func ExpensesServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Expense 1")
}
