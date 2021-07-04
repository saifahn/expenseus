package main

import (
	"log"
	"net/http"

	"github.com/saifahn/expenseus"
)

func main() {
	webservice := &expenseus.WebService{}
	log.Fatal(http.ListenAndServe(":5000", webservice))
}
