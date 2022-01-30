package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/securecookie"
	"github.com/saifahn/expenseus/internal/app"
	"github.com/saifahn/expenseus/internal/ddb"
	"github.com/saifahn/expenseus/internal/googleoauth"
	"github.com/saifahn/expenseus/internal/s3images"
	"github.com/saifahn/expenseus/internal/sessions"
)

func main() {
	var (
		frontendURL       string
		useLocalAWSConfig bool
	)

	if mode := os.Getenv("MODE"); mode == "development" {
		frontendURL = os.Getenv("FRONTEND_DEV_SERVER")
		useLocalAWSConfig = true
	} else {
		frontendURL = "/"
		useLocalAWSConfig = false
	}

	var dynamo *dynamodb.DynamoDB

	if useLocalAWSConfig {
		sess := session.Must(session.NewSession(aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(os.Getenv("DYNAMODB_DEV_ID"), os.Getenv("DYNAMODB_DEV_SECRET"), ""))))
		sess.Config.Endpoint = aws.String(os.Getenv("DYNAMODB_ENDPOINT_LOCAL"))
		dynamo = dynamodb.New(sess)
	}

	db := ddb.New(dynamo, os.Getenv("DYNAMODB_USERS_TABLE_NAME"), os.Getenv("DYNAMODB_TRANSACTIONS_TABLE_NAME"))

	googleOauth := googleoauth.New()

	tempHashKey := securecookie.GenerateRandomKey(64)
	tempBlockKey := securecookie.GenerateRandomKey(32)

	sessions := sessions.New(tempHashKey, tempBlockKey)
	images := s3images.New(useLocalAWSConfig)

	wb := app.NewWebService(db, googleOauth, sessions, frontendURL, images)

	r := app.InitRouter(wb)

	log.Fatal(http.ListenAndServe(":5000", r))
}
