package main

import (
	"log"
	"net/http"

	"github.com/saifahn/expenseus"
	"github.com/saifahn/expenseus/redis"
)

var redisAddr = "localhost:6379"

func main() {
	rdb := redis.New(redisAddr)

	wb := expenseus.NewWebService(rdb)

	r := expenseus.InitRouter(wb)

	log.Fatal(http.ListenAndServe(":5000", r))
}
