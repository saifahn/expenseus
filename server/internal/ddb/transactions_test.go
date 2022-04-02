package ddb

import (
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/stretchr/testify/assert"
)

const testTransactionsTableName = "expenseus-testing-transactions"

func TestTransactionTable(t *testing.T) {
	assert := assert.New(t)
	dynamodb := NewDynamoDBLocalAPI()

	// create the table in the local test database
	err := CreateTestTable(dynamodb, testTransactionsTableName)
	if err != nil {
		t.Logf("table could not be created: %v", err)
	}
	defer DeleteTable(dynamodb, testTransactionsTableName)
	tbl := table.New(dynamodb, testTransactionsTableName)
	// create the transactions table instance
	transactions := NewTransactionsTable(tbl)

	// retrieving a non-existent item will give an error
	_, err = transactions.Get("non-existent-user", "non-existent-item")
	assert.EqualError(err, table.ErrItemNotFound.Error())

	item := &TransactionItem{
		PK:         "user#test-user-id",
		SK:         "txn#test-txn-id",
		ID:         "test-txn-id",
		UserID:     "test-user-id",
		EntityType: "transaction",
		GSI1PK:     "transactions",
		GSI1SK:     "txn#test-txn-id",
	}

	// no error raised the first time
	err = transactions.PutIfNotExists(*item)
	assert.NoError(err)

	// it is possible to overwrite with Put
	err = transactions.Put(*item)
	assert.NoError(err)

	// it will now raise an error as the item exists
	err = transactions.PutIfNotExists(*item)
	assert.EqualError(err, ErrConflict.Error())

	// the item is successfully retrieved
	got, err := transactions.Get(item.UserID, item.ID)
	assert.NoError(err)
	assert.Equal(item, got)

	// get all transactions
	itemsGot, err := transactions.GetAll()
	assert.NoError(err)
	assert.Len(itemsGot, 1)
	assert.Contains(itemsGot, *item)

	// get the transactions by username
	itemsGot, err = transactions.GetByUserID(item.UserID)
	assert.NoError(err)
	assert.Contains(itemsGot, *item)

	// the item is successfully deleted
	err = transactions.Delete(item.UserID, item.ID)
	assert.NoError(err)
	_, err = transactions.Get(item.UserID, item.ID)
	assert.EqualError(err, table.ErrItemNotFound.Error())
}
