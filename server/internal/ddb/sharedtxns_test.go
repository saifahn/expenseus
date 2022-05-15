package ddb

import (
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/stretchr/testify/assert"
)

const testSharedTxnsTableName = "expenseus-testing-shared-txns"

func TestCreateSharedTxn(t *testing.T) {
	assert := assert.New(t)
	dynamoDB := NewDynamoDBLocalAPI()
	err := CreateTestTable(dynamoDB, testSharedTxnsTableName)
	if err != nil {
		t.Fatalf("table could not be created: %v", err)
	}
	defer DeleteTable(dynamoDB, testSharedTxnsTableName)
	tbl := table.New(dynamoDB, testSharedTxnsTableName)
	sharedTxns := NewSharedTxnsRepository(tbl)

	testUserIDs := []string{"test-01", "test-01"}
	testInput := CreateSharedTxnInput{
		ID:           "test-shared-txn-id",
		TrackerID:    " test-tracker-id",
		Participants: testUserIDs,
	}

	err = sharedTxns.Create(testInput)
	assert.NoError(err)
}
