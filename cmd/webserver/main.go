package main

import (
	"log"
	"net/http"

	"github.com/saifahn/expenseus"
	"github.com/saifahn/expenseus/redis"
)

func main() {
	rdb := redis.New()

	wb := expenseus.NewWebService(rdb)

	r := expenseus.InitRouter(wb)

	log.Fatal(http.ListenAndServe(":5000", r))
}
