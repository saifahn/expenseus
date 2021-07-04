package main

import (
	"log"
	"net/http"

	"github.com/saifahn/expenseus"
)

func main() {
	handler := http.HandlerFunc(expenseus.WebService)
	log.Fatal(http.ListenAndServe(":5000", handler))
}
