package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/saifahn/expenseus"
	"github.com/saifahn/expenseus/googleoauth"
	"github.com/saifahn/expenseus/redis"
	"github.com/saifahn/expenseus/sessions"
)

var redisAddr = os.Getenv("REDIS_ADDRESS")

func main() {
	var frontendURL string

	if mode := os.Getenv("MODE"); mode == "development" {
		frontendURL = os.Getenv("FRONTEND_DEV_SERVER")
	} else {
		frontendURL = "/"
	}

	rdb := redis.New(redisAddr)

	googleOauth := googleoauth.New()

	tempHashKey := securecookie.GenerateRandomKey(64)
	tempBlockKey := securecookie.GenerateRandomKey(32)

	sessions := sessions.New(tempHashKey, tempBlockKey)

	wb := expenseus.NewWebService(rdb, googleOauth, sessions, frontendURL)

	r := expenseus.InitRouter(wb)

	log.Fatal(http.ListenAndServe(":5000", r))
}
