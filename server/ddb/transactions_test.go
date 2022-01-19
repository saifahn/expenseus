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
	tbl := table.New(dynamodb, testTransactionsTableName)
	// create the transactions table instance
	transactions := NewTransactionsTable(tbl)

	// retrieving a non-existent item will give an error
	_, err = transactions.Get("non-existent-item")
	assert.EqualError(err, table.ErrItemNotFound.Error())

	item := &TransactionItem{
		ID: "test-item-id",
		// Amount: 123,
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
	got, err := transactions.Get(item.ID)
	assert.NoError(err)
	assert.Equal(item, got)

	// the item is successfully deleted
	err = transactions.Delete(item.ID)
	assert.NoError(err)
	_, err = transactions.Get(item.ID)
	assert.EqualError(err, table.ErrItemNotFound.Error())
	DeleteTable(dynamodb, testTransactionsTableName)
}

func TestGetAll(t *testing.T) {
	assert := assert.New(t)
	dynamodb := NewDynamoDBLocalAPI()
	err := CreateTestTable(dynamodb, testTransactionsTableName)
	if err != nil {
		t.Logf("table could not be created: %v", err)
	}

	tbl := table.New(dynamodb, testTransactionsTableName)
	// create the transactions table instance
	transactions := NewTransactionsTable(tbl)

	item := TransactionItem{
		ID: "test-item-id",
	}

	// no error raised the first time
	err = transactions.PutIfNotExists(item)
	assert.NoError(err)

	items, err := transactions.GetAll()
	assert.NoError(err)
	assert.Len(items, 1)
	assert.Contains(items, item)
	DeleteTable(dynamodb, testTransactionsTableName)
}
