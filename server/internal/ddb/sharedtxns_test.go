package ddb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSharedTxnsTableName = "expenseus-testing-shared-txns"

func TestSharedTxns(t *testing.T) {
	assert := assert.New(t)
	tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
	defer teardown()
	sharedTxns := NewSharedTxnsRepository(tbl)

	testUserIDs := []string{"test-01", "test-01"}
	testInput := CreateSharedTxnInput{
		ID:           "test-shared-txn-id",
		TrackerID:    "test-tracker-id",
		Participants: testUserIDs,
	}

	err := sharedTxns.Create(testInput)
	assert.NoError(err)

	got, err := sharedTxns.GetFromTracker(testInput.TrackerID)
	assert.NoError(err)

	want := SharedTxnItem{
		PK:           makeTrackerIDKey(testInput.TrackerID),
		SK:           makeSharedTxnIDKey(testInput.ID),
		EntityType:   sharedTxnEntityType,
		ID:           testInput.ID,
		Tracker:      testInput.TrackerID,
		Participants: testInput.Participants,
	}
	assert.Contains(got, want)
}
