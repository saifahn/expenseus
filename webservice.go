package expenseus

import (
	"fmt"
	"net/http"
)

func WebService(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Expense 1")
}
