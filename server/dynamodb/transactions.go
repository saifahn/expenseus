package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
	"github.com/pkg/errors"
)

type TransactionItem struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
}

type TransactionsTable interface {
	Get(id string) (*TransactionItem, error)
	PutIfNotExists(item TransactionItem) error
	Delete(id string) error
}

type transactionsTable struct {
	table *table.Table
}

const TransactionsHashKeyName = "id"

func conflictOrErr(err error) error {
	dynamoErr, ok := errors.Cause(err).(awserr.Error)
	if ok && dynamoErr.Code() == "ConditionalCheckFailedException" {
		return errors.New("dynamodb: conflict")
	}
	return err
}

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

func (t *transactionsTable) Delete(id string) error {
	return t.table.DeleteItem(attributes.String(id), nil)
}
