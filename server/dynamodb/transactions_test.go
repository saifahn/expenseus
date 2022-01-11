package dynamodb

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/stretchr/testify/assert"
)

func newDynamoDBLocalAPI() dynamodbiface.DynamoDBAPI {
	sess := session.Must(
		session.NewSession(aws.NewConfig().WithCredentials(credentials.NewStaticCredentials("dynamodb-testing", "test-secret", ""))),
	)
	// TODO: replace with environment variables?
	sess.Config.Endpoint = aws.String("http://localhost:8000")
	// sess.Config.Region = aws.String("ap-")
	return dynamodb.New(sess)
}

const testTableName = "expenseus-testing-transactions"

func createTestTable(d dynamodbiface.DynamoDBAPI) error {
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
		TableName: aws.String(testTableName),
	}
	_, err := d.CreateTable(input)
	if err != nil {
		return err
	}

	fmt.Println("successfully created the table", testTableName)
	return nil
}

func TestTransactionTable(t *testing.T) {
	assert := assert.New(t)
	dynamodb := newDynamoDBLocalAPI()

	// create the table in the local test database
	err := createTestTable(dynamodb)
	if err != nil {
		t.Logf("table could not be crated: %v", err)
	}
	tbl := table.New(dynamodb, testTableName)
	// create the transactions table
	transactions := NewTransactionsTable(tbl)

	_, err = transactions.Get("non-existent-item")
	assert.Equal(table.ErrItemNotFound, err)

	item := &TransactionItem{
		ID:     "test-item-id",
		Amount: 123,
	}

	err = transactions.PutIfNotExists(*item)
	assert.NoError(err)
}
