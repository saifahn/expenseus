package ddb

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// NewDynamoDBLocalAPI creates a new session with DynamoDB for local testing.
func NewDynamoDBLocalAPI() dynamodbiface.DynamoDBAPI {
	sess := session.Must(
		session.NewSession(aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(os.Getenv("DYNAMODB_TEST_ID"), os.Getenv("DYNAMODB_TEST_SECRET"), ""))),
	)
	sess.Config.Endpoint = aws.String(os.Getenv("DYNAMODB_ENDPOINT_LOCAL"))
	return dynamodb.New(sess)
}

// CreateTestTable creates a table for testing.
func CreateTestTable(d dynamodbiface.DynamoDBAPI, name string) error {
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
			{
				AttributeName: aws.String("GSI1PK"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("GSI1SK"),
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
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{{
			IndexName: aws.String("GSI1"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("GSI1PK"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("GSI1SK"),
					KeyType:       aws.String("RANGE"),
				},
			},
			Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		}},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(name),
	})
	if err != nil {
		return err
	}

	log.Println("successfully created the table", name)
	return nil
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
