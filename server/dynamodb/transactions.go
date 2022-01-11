package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
)

type TransactionItem struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
}

type TransactionsTable interface {
	Get(id string) (*TransactionItem, error)
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
