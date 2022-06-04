package ddb

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	tablePrimaryKey        = "PK"
	tableSortKey           = "SK"
	gsi1Name               = "GSI1"
	gsi1PrimaryKey         = "GSI1PK"
	gsi1SortKey            = "GSI1SK"
	unsettledTxnsIndexName = "UnsettledTransactions"
	unsettledTxnsIndexPK   = "PK"
	unsettledTxnsIndexSK   = "Unsettled"
)

// CreateTable creates a table of the expenseus schema with the given name.
func CreateTable(d dynamodbiface.DynamoDBAPI, name string) error {
	_, err := d.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(tablePrimaryKey),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(tableSortKey),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(gsi1PrimaryKey),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(gsi1SortKey),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(unsettledTxnsIndexSK),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(tablePrimaryKey),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String(tableSortKey),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String(gsi1Name),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String(gsi1PrimaryKey),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String(gsi1SortKey),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
			{
				IndexName: aws.String(unsettledTxnsIndexName),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String(unsettledTxnsIndexPK),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String(unsettledTxnsIndexSK),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
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

	log.Println("successfully created the table", name)
	return nil
}
