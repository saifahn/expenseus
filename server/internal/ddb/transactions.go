package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
)

const (
	txnKeyPrefix  = "txn"
	txnEntityType = "transaction"
	allTxnKey     = "transactions"
)

type TransactionItem struct {
	PK         string `json:"PK"`
	SK         string `json:"SK"`
	EntityType string `json:"EntityType"`
	ID         string `json:"ID"`
	UserID     string `json:"UserID"`
	Name       string `json:"Name"`
	Amount     int64  `json:"Amount"`
	Date       int64  `json:"Date"`
	GSI1PK     string `json:"GSI1PK"`
	GSI1SK     string `json:"GSI1SK"`
}

type GetTxnInput struct {
	ID     string
	UserID string
}

type TxnRepository interface {
	Get(transactionID string) (*TransactionItem, error)
	GetOne(input GetTxnInput) (*TransactionItem, error)
	GetAll() ([]TransactionItem, error)
	GetByUserID(id string) ([]TransactionItem, error)
	Create(item TransactionItem) error
	Update(item TransactionItem) error
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

func (t *txnRepo) Get(txnID string) (*TransactionItem, error) {
	txnIDKey := makeTxnIDKey(txnID)

	options := []option.QueryInput{
		option.Index("GSI1"),
		option.QueryExpressionAttributeName(gsi1PrimaryKey, "#GSI1PK"),
		option.QueryExpressionAttributeName(gsi1SortKey, "#GSI1SK"),
		option.QueryExpressionAttributeValue(":allTransactionsKey", attributes.String(allTxnKey)),
		option.QueryExpressionAttributeValue(":transactionID", attributes.String(txnIDKey)),
		option.QueryKeyConditionExpression("#GSI1PK = :allTransactionsKey AND #GSI1SK = :transactionID"),
	}

	var items []TransactionItem

	_, err := t.table.Query(&items, options...)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, table.ErrItemNotFound
	}
	if len(items) > 1 {
		return nil, ErrUnexpected
	}

	return &items[0], nil
}

func (t *txnRepo) GetOne(input GetTxnInput) (*TransactionItem, error) {
	userIDKey := makeUserIDKey(input.UserID)
	txnIDKey := makeTxnIDKey(input.ID)

	item := &TransactionItem{}
	err := t.table.GetItem(attributes.String(userIDKey), attributes.String(txnIDKey), item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *txnRepo) Create(item TransactionItem) error {
	err := t.table.PutItem(
		item,
		option.PutCondition("attribute_not_exists(SK)"),
	)
	if err != nil {
		return conflictOrErr(err)
	}

	return nil
}

func (t *txnRepo) Update(item TransactionItem) error {
	// this condition is a stand in for the item existing, as the primary key
	// and index keys must be present for a item to exist
	err := t.table.PutItem(item, option.PutCondition("attribute_exists(SK)"))
	if err != nil {
		return attrNotExistsOrErr(err)
	}

	return nil
}

func (t *txnRepo) Delete(txnID, userID string) error {
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

func (t *txnRepo) GetByUserID(id string) ([]TransactionItem, error) {
	userIDKey := makeUserIDKey(id)
	allTxnPrefix := fmt.Sprintf("%s#", txnKeyPrefix)

	options := []option.QueryInput{
		option.QueryExpressionAttributeName(tablePrimaryKey, "#PK"),
		option.QueryExpressionAttributeName(tableSortKey, "#SK"),
		option.QueryExpressionAttributeValue(":userKey", attributes.String(userIDKey)),
		option.QueryExpressionAttributeValue(":allTxnPrefix", attributes.String(allTxnPrefix)),
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
