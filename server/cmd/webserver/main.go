package main

import (
	"log"
	"net/http"

	"github.com/saifahn/expenseus"
	"github.com/saifahn/expenseus/googleoauth"
	"github.com/saifahn/expenseus/redis"
)

var redisAddr = "localhost:6379"

func main() {
	rdb := redis.New(redisAddr)

	googleOauth := googleoauth.New()

	wb := expenseus.NewWebService(rdb, googleOauth)

	r := expenseus.InitRouter(wb)

	log.Fatal(http.ListenAndServe(":5000", r))
}
