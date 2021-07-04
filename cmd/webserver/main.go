package main

import (
	"log"
	"net/http"

	"github.com/saifahn/expenseus"
)

type InMemoryExpenseStore struct{}

func (i *InMemoryExpenseStore) GetExpense(id string) string {
	return "123"
}

func main() {
	webservice := expenseus.NewWebService(&InMemoryExpenseStore{})
	log.Fatal(http.ListenAndServe(":5000", webservice))
}
