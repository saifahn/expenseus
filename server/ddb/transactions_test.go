package ddb

import (
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/saifahn/expenseus"
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
	_, err = transactions.Get("non-existent-item")
	assert.EqualError(err, table.ErrItemNotFound.Error())

	testED := &expenseus.ExpenseDetails{
		UserID: "test-user",
		Name:   "test-expense",
	}

	item := &TransactionItem{
		ID:             "test-item-id",
		ExpenseDetails: *testED,
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

	// get all expenses
	itemsGot, err := transactions.GetAll()
	assert.NoError(err)
	assert.Len(itemsGot, 1)
	assert.Contains(itemsGot, *item)

	// get the expenses by username
	itemsGot, err = transactions.GetByUserID(testED.UserID)
	assert.NoError(err)
	assert.Contains(itemsGot, *item)

	// the item is successfully deleted
	err = transactions.Delete(item.ID)
	assert.NoError(err)
	_, err = transactions.Get(item.ID)
	assert.EqualError(err, table.ErrItemNotFound.Error())
}
