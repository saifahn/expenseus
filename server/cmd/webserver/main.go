package main

import (
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/saifahn/expenseus"
	"github.com/saifahn/expenseus/googleoauth"
	"github.com/saifahn/expenseus/redis"
	"github.com/saifahn/expenseus/sessions"
)

var redisAddr = "localhost:6379"

func main() {
	rdb := redis.New(redisAddr)

	googleOauth := googleoauth.New()

	tempHashKey := securecookie.GenerateRandomKey(64)
	tempBlockKey := securecookie.GenerateRandomKey(32)

	sessions := sessions.New(tempHashKey, tempBlockKey)

	wb := expenseus.NewWebService(rdb, googleOauth, sessions)

	r := expenseus.InitRouter(wb)

	log.Fatal(http.ListenAndServe(":5000", r))
}
