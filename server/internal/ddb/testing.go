package ddb

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nabeken/aws-go-dynamodb/table"
)

// NewDynamoDBLocalAPI creates a new session with DynamoDB for local testing.
func NewDynamoDBLocalAPI() dynamodbiface.DynamoDBAPI {
	sess := session.Must(
		session.NewSession(aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(os.Getenv("DYNAMODB_TEST_ID"), os.Getenv("DYNAMODB_TEST_SECRET"), ""))),
	)
	sess.Config.Endpoint = aws.String(os.Getenv("DYNAMODB_ENDPOINT_LOCAL"))
	return dynamodb.New(sess)
}

// DeleteTable deletes a table. It should only be used in testing.
func DeleteTable(d dynamodbiface.DynamoDBAPI, name string) error {
	_, err := d.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(name)})
	if err != nil {
		return err
	}

	log.Println("successfully deleted the table", name)
	return nil
}

// SetUpTestTable creates a table for testing.
func SetUpTestTable(t testing.TB, tableName string) (*table.Table, func()) {
	ddb := NewDynamoDBLocalAPI()
	err := CreateTable(ddb, tableName)
	if err != nil {
		t.Fatalf("table could not be created: %v", err)
	}

	teardown := func() {
		DeleteTable(ddb, tableName)
	}
	return table.New(ddb, tableName), teardown
}
