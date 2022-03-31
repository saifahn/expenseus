package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
)

type TransactionItem struct {
	PK         string `json:"PK"`
	SK         string `json:"SK"`
	ID         string `json:"ID"`
	UserID     string `json:"UserID"`
	EntityType string `json:"EntityType"`
}

const TransactionKeyPrefix = "txn"

type TransactionsTable interface {
	Get(userID, transactionID string) (*TransactionItem, error)
	// GetAll() ([]TransactionItem, error)
	// GetByUserID(userID string) ([]TransactionItem, error)
	PutIfNotExists(item TransactionItem) error
	Put(item TransactionItem) error
	Delete(userID, transactionID string) error
}

type transactionsTable struct {
	table *table.Table
}

func NewTransactionsTable(t *table.Table) TransactionsTable {
	t.WithHashKey(HashKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(RangeKey, dynamodb.ScalarAttributeTypeS)
	return &transactionsTable{table: t}
}

func (t *transactionsTable) Get(userID, txnID string) (*TransactionItem, error) {
	userKey := fmt.Sprintf("%s#%s", UserKeyPrefix, userID)
	txnKey := fmt.Sprintf("%s#%s", TransactionKeyPrefix, txnID)
	item := &TransactionItem{}
	err := t.table.GetItem(attributes.String(userKey), attributes.String(txnKey), item, option.ConsistentRead())
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *transactionsTable) PutIfNotExists(item TransactionItem) error {
	err := t.table.PutItem(
		item,
		option.PutCondition("attribute_not_exists(SK)"),
	)
	if err != nil {
		return conflictOrErr(err)
	}

	return nil
}

func (t *transactionsTable) Put(item TransactionItem) error {
	return t.table.PutItem(item)
}

func (t *transactionsTable) Delete(userID, txnID string) error {
	userKey := fmt.Sprintf("%s#%s", UserKeyPrefix, userID)
	txnKey := fmt.Sprintf("%s#%s", TransactionKeyPrefix, txnID)
	return t.table.DeleteItem(attributes.String(userKey), attributes.String(txnKey))
}

// func (t *transactionsTable) GetAll() ([]TransactionItem, error) {
// 	response, err := t.table.DynamoDB.Scan(&dynamodb.ScanInput{TableName: t.table.Name})
// 	if err != nil {
// 		return nil, err
// 	}
// 	var items []TransactionItem

// 	for _, i := range response.Items {
// 		var item TransactionItem
// 		err = dynamodbattribute.UnmarshalMap(i, &item)
// 		if err != nil {
// 			return nil, err
// 		}
// 		items = append(items, item)
// 	}

// 	return items, nil
// }

// func (t *transactionsTable) GetByUserID(userid string) ([]TransactionItem, error) {
// 	filt := expression.Name("userId").Equal(expression.Value(userid))
// 	expr, err := expression.NewBuilder().WithFilter(filt).Build()
// 	if err != nil {
// 		return nil, err
// 	}

// 	params := &dynamodb.ScanInput{
// 		ExpressionAttributeNames:  expr.Names(),
// 		ExpressionAttributeValues: expr.Values(),
// 		FilterExpression:          expr.Filter(),
// 		ProjectionExpression:      expr.Projection(),
// 		TableName:                 t.table.Name,
// 	}

// 	result, err := t.table.DynamoDB.Scan(params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var items []TransactionItem

// 	for _, i := range result.Items {
// 		var item TransactionItem
// 		err = dynamodbattribute.UnmarshalMap(i, &item)
// 		if err != nil {
// 			return nil, err
// 		}
// 		items = append(items, item)
// 	}

// 	return items, nil
// }
