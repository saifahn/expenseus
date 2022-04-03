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
	txnKeyPrefix  = "txn"
	txnEntityType = "transaction"
	allTxnKey     = "transactions"
)

type TxnRepository interface {
	Get(userID, transactionID string) (*TransactionItem, error)
	GetAll() ([]TransactionItem, error)
	GetByUserID(userID string) ([]TransactionItem, error)
	PutIfNotExists(item TransactionItem) error
	Put(item TransactionItem) error
	Delete(userID, transactionID string) error
}

type txnRepo struct {
	table *table.Table
}

func NewTxnRepository(t *table.Table) TxnRepository {
	t.WithHashKey(tablePrimaryKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(tableSortKey, dynamodb.ScalarAttributeTypeS)
	return &txnRepo{table: t}
}

func (t *txnRepo) Get(userID, txnID string) (*TransactionItem, error) {
	userIDKey := makeUserIDKey(userID)
	txnIDKey := makeTxnIDKey(txnID)
	item := &TransactionItem{}
	err := t.table.GetItem(attributes.String(userIDKey), attributes.String(txnIDKey), item, option.ConsistentRead())
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *txnRepo) PutIfNotExists(item TransactionItem) error {
	err := t.table.PutItem(
		item,
		option.PutCondition("attribute_not_exists(SK)"),
	)
	if err != nil {
		return conflictOrErr(err)
	}

	return nil
}

func (t *txnRepo) Put(item TransactionItem) error {
	return t.table.PutItem(item)
}

func (t *txnRepo) Delete(userID, txnID string) error {
	userIDKey := makeUserIDKey(userID)
	txnIDKey := makeTxnIDKey(txnID)
	return t.table.DeleteItem(attributes.String(userIDKey), attributes.String(txnIDKey))
}

func (t *txnRepo) GetAll() ([]TransactionItem, error) {
	options := []option.QueryInput{
		option.Index("GSI1"),
		option.QueryExpressionAttributeName(gsi1PrimaryKey, "#GSI1PK"),
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

func (t *txnRepo) GetByUserID(userID string) ([]TransactionItem, error) {
	userIDKey := makeUserIDKey(userID)
	// to match all transactions
	txnKeyWithoutID := fmt.Sprintf("%s#", txnKeyPrefix)

	options := []option.QueryInput{
		option.QueryExpressionAttributeName(tablePrimaryKey, "#PK"),
		option.QueryExpressionAttributeName(tableSortKey, "#SK"),
		option.QueryExpressionAttributeValue(":userKey", attributes.String(userIDKey)),
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

func makeTxnIDKey(id string) string {
	return fmt.Sprintf("%s#%s", txnKeyPrefix, id)
}
