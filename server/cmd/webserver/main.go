package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/securecookie"
	"github.com/saifahn/expenseus/internal/app"
	"github.com/saifahn/expenseus/internal/ddb"
	"github.com/saifahn/expenseus/internal/googleoauth"
	"github.com/saifahn/expenseus/internal/router"
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

	err := ddb.CreateTable(dynamo, os.Getenv("DYNAMODB_TABLE_NAME"))
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceInUseException {
				log.Print("table already exists")
			} else {
				log.Print(err.Error())
			}
		} else {
			log.Print(err.Error())
		}
	}

	db := ddb.New(dynamo, os.Getenv("DYNAMODB_TABLE_NAME"))

	googleOauth := googleoauth.New()

	tempHashKey := securecookie.GenerateRandomKey(64)
	tempBlockKey := securecookie.GenerateRandomKey(32)

	sessions := sessions.New(tempHashKey, tempBlockKey)
	images := s3images.New(useLocalAWSConfig)

	a := app.New(db, googleOauth, sessions, frontendURL, images)

	r := router.Init(a)

	log.Fatal(http.ListenAndServe(":4000", r))
}
