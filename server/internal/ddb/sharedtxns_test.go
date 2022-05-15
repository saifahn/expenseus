package ddb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSharedTxnsTableName = "expenseus-testing-shared-txns"

func TestCreateSharedTxn(t *testing.T) {
	assert := assert.New(t)
	tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
	defer teardown()
	sharedTxns := NewSharedTxnsRepository(tbl)

	testUserIDs := []string{"test-01", "test-01"}
	testInput := CreateSharedTxnInput{
		ID:           "test-shared-txn-id",
		TrackerID:    " test-tracker-id",
		Participants: testUserIDs,
	}

	err := sharedTxns.Create(testInput)
	assert.NoError(err)
}
