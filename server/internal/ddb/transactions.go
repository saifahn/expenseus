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
	EntityType string `json:"EntityType"`
	ID         string `json:"ID"`
	UserID     string `json:"UserID"`
	GSI1PK     string `json:"GSI1PK"`
	GSI1SK     string `json:"GSI1SK"`
}

const (
	TransactionKeyPrefix  = "txn"
	transactionEntityType = "transaction"
	allTxnKey             = "transactions"
)

type TransactionsTable interface {
	Get(userID, transactionID string) (*TransactionItem, error)
	GetAll() ([]TransactionItem, error)
	GetByUserID(userID string) ([]TransactionItem, error)
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

func (t *transactionsTable) GetAll() ([]TransactionItem, error) {
	options := []option.QueryInput{
		option.Index("GSI1"),
		option.QueryExpressionAttributeName(GSI1PK, "#GSI1PK"),
		option.QueryExpressionAttributeValue(":allTransactionsKey", attributes.String(allTxnKey)),
		option.QueryKeyConditionExpression("#GSI1PK = :allTransactionsKey"),
	}

	var items []TransactionItem

	_, err := t.table.Query(&items, options...)

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (t *transactionsTable) GetByUserID(userID string) ([]TransactionItem, error) {
	userKey := fmt.Sprintf("%s#%s", UserKeyPrefix, userID)
	txnKeyWithoutID := fmt.Sprintf("%s#", TransactionKeyPrefix)

	options := []option.QueryInput{
		option.QueryExpressionAttributeName(HashKey, "#PK"),
		option.QueryExpressionAttributeName(RangeKey, "#SK"),
		option.QueryExpressionAttributeValue(":userKey", attributes.String(userKey)),
		option.QueryExpressionAttributeValue(":allTxnPrefix", attributes.String(txnKeyWithoutID)),
		option.QueryKeyConditionExpression("#PK = :userKey and begins_with(#SK, :allTxnPrefix)"),
	}

	var items []TransactionItem

	_, err := t.table.Query(&items, options...)

	if err != nil {
		return nil, err
	}

	return items, nil
}
