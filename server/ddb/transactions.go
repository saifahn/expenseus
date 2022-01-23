package ddb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
	"github.com/saifahn/expenseus"
)

type TransactionItem struct {
	expenseus.ExpenseDetails
	ID string `json:"id"`
}

type TransactionsTable interface {
	Get(id string) (*TransactionItem, error)
	GetAll() ([]TransactionItem, error)
	GetByUserID(userid string) ([]TransactionItem, error)
	PutIfNotExists(item TransactionItem) error
	Put(item TransactionItem) error
	Delete(id string) error
}

type transactionsTable struct {
	table *table.Table
}

const TransactionsHashKeyName = "id"

func NewTransactionsTable(t *table.Table) TransactionsTable {
	t.WithHashKey(TransactionsHashKeyName, dynamodb.ScalarAttributeTypeS)
	return &transactionsTable{table: t}
}

func (t *transactionsTable) Get(id string) (*TransactionItem, error) {
	item := &TransactionItem{}
	err := t.table.GetItem(attributes.String(id), nil, item, option.ConsistentRead())
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *transactionsTable) PutIfNotExists(item TransactionItem) error {
	err := t.table.PutItem(
		item,
		option.PutCondition("attribute_not_exists(id)"),
	)
	if err != nil {
		return conflictOrErr(err)
	}

	return nil
}

func (t *transactionsTable) Put(item TransactionItem) error {
	return t.table.PutItem(item)
}

func (t *transactionsTable) Delete(id string) error {
	return t.table.DeleteItem(attributes.String(id), nil)
}

func (t *transactionsTable) GetAll() ([]TransactionItem, error) {
	response, err := t.table.DynamoDB.Scan(&dynamodb.ScanInput{TableName: t.table.Name})
	if err != nil {
		return nil, err
	}
	var items []TransactionItem

	for _, i := range response.Items {
		var item TransactionItem
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (t *transactionsTable) GetByUserID(userid string) ([]TransactionItem, error) {
	filt := expression.Name("userId").Equal(expression.Value(userid))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 t.table.Name,
	}

	result, err := t.table.DynamoDB.Scan(params)
	if err != nil {
		return nil, err
	}

	var items []TransactionItem

	for _, i := range result.Items {
		var item TransactionItem
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
