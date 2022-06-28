package ddb

import (
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/stretchr/testify/assert"
)

const testTransactionsTableName = "expenseus-testing-transactions"

func TestTransactionTable(t *testing.T) {
	assert := assert.New(t)
	tbl, teardown := SetUpTestTable(t, testTransactionsTableName)
	defer teardown()
	transactions := NewTxnRepository(tbl)

	// retrieving a non-existent item will give an error
	_, err := transactions.Get(GetTxnInput{
		UserID: "non-existent-user-id",
		ID:     "non-existent-item-id",
	})
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
	err = transactions.Create(*item)
	assert.NoError(err)

	// it will now raise an error as the item exists
	err = transactions.Create(*item)
	assert.EqualError(err, ErrConflict.Error())

	// the item is successfully retrieved
	got, err := transactions.Get(GetTxnInput{
		UserID: item.UserID,
		ID:     item.ID,
	})
	assert.NoError(err)
	assert.Equal(item, got)

	// get the transactions by username
	itemsGot, err := transactions.GetByUserID(item.UserID)
	assert.NoError(err)
	assert.Contains(itemsGot, *item)

	// the item is successfully deleted
	err = transactions.Delete(item.ID, item.UserID)
	assert.NoError(err)
	_, err = transactions.Get(GetTxnInput{
		UserID: item.UserID,
		ID:     item.ID,
	})
	assert.ErrorIs(err, table.ErrItemNotFound)
}

func TestUpdateItem(t *testing.T) {
	initialItem := &TransactionItem{
		PK:       "user#test-user-id",
		SK:       "txn#test-txn-id",
		GSI1PK:   "transactions",
		GSI1SK:   "txn#test-txn-id",
		ID:       "test-txn-id",
		UserID:   "test-user-id",
		Location: "initial-location",
		Details:  "original-transaction",
	}
	updatedItem := &TransactionItem{
		PK:       initialItem.PK,
		SK:       initialItem.SK,
		GSI1PK:   initialItem.GSI1PK,
		GSI1SK:   initialItem.GSI1SK,
		ID:       initialItem.ID,
		UserID:   initialItem.UserID,
		Location: "changed-location",
		Details:  "transaction-name-changed",
	}

	tests := map[string]struct {
		initialItem  *TransactionItem
		itemToUpdate *TransactionItem
		finalItem    *TransactionItem
		wantErr      error
	}{
		"updating a non-existent item will give an error": {
			initialItem: initialItem,
			itemToUpdate: &TransactionItem{
				PK:       "user#a-different-user",
				SK:       "txn#a-different-item",
				GSI1PK:   "transactions",
				GSI1SK:   "txn#a-different-item",
				ID:       "a-different-item",
				UserID:   "different-user-id",
				Location: "not-the-original",
			},
			finalItem: initialItem,
			wantErr:   ErrAttrNotExists,
		},
		"updating an existing item will update it as expected": {
			initialItem:  initialItem,
			itemToUpdate: updatedItem,
			finalItem:    updatedItem,
			wantErr:      nil,
		},
	}

	for name, tc := range tests {
		assert := assert.New(t)

		t.Run(name, func(t *testing.T) {
			tbl, teardown := SetUpTestTable(t, "test-txn-update-items")
			defer teardown()
			transactions := NewTxnRepository(tbl)

			err := transactions.Create(*tc.initialItem)
			assert.NoError(err)

			err = transactions.Update(*tc.itemToUpdate)
			if tc.wantErr != nil {
				assert.ErrorIs(err, ErrAttrNotExists)
			}

			got, err := transactions.Get(
				GetTxnInput{
					UserID: tc.initialItem.UserID,
					ID:     tc.initialItem.ID,
				},
			)
			assert.NoError(err)
			assert.Equal(tc.finalItem, got)
		})
	}
}

func TestGetBetweenDates(t *testing.T) {
	initialItem := TransactionItem{
		PK:       "user#test-user-id",
		SK:       "txn#test-txn-id",
		GSI1PK:   "user#test-user-id",
		GSI1SK:   "txn#10000#test-txn-id",
		ID:       "test-txn-id",
		UserID:   "test-user-id",
		Location: "initial-location",
		Date:     10000,
	}

	tests := map[string]struct {
		wantItems []TransactionItem
		from      int64
		to        int64
	}{
		"with a date-range containing an item": {
			wantItems: []TransactionItem{initialItem},
			from:      10000,
			to:        20000,
		},
		"with a date range not containing any items": {
			wantItems: []TransactionItem{},
			from:      15000,
			to:        20000,
		},
	}

	for name, tc := range tests {
		assert := assert.New(t)

		t.Run(name, func(t *testing.T) {
			tbl, teardown := SetUpTestTable(t, "test-get-txns-between-dates")
			defer teardown()
			txns := NewTxnRepository(tbl)

			err := txns.Create(initialItem)
			assert.NoError(err)

			got, err := txns.GetBetweenDates(initialItem.UserID, tc.from, tc.to)
			assert.NoError(err)
			assert.ElementsMatch(got, tc.wantItems)
		})
	}
}
