package ddb

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func NewDynamoDBLocalAPI() dynamodbiface.DynamoDBAPI {
	sess := session.Must(
		session.NewSession(aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(os.Getenv("DYNAMODB_TEST_ID"), os.Getenv("DYNAMODB_TEST_SECRET"), ""))),
	)
	sess.Config.Endpoint = aws.String(os.Getenv("DYNAMODB_ENDPOINT_LOCAL"))
	return dynamodb.New(sess)
}

func CreateExpenseusTable(d dynamodbiface.DynamoDBAPI, name string) error {
	_, err := d.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(name),
	})
	if err != nil {
		return err
	}

	fmt.Println("successfully created the table", name)
	return nil
}

func CreateTestTable(d dynamodbiface.DynamoDBAPI, name string) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(name),
	}
	_, err := d.CreateTable(input)
	if err != nil {
		return err
	}

	fmt.Println("successfully created the table", name)
	return nil
}

func DeleteTable(d dynamodbiface.DynamoDBAPI, name string) error {
	_, err := d.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(name)})
	if err != nil {
		return err
	}

	fmt.Println("successfully deleted the table", name)
	return nil
}
