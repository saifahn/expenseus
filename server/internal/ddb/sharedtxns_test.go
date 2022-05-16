package ddb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSharedTxnsTableName = "expenseus-testing-shared-txns"

func TestSharedTxns(t *testing.T) {
	assert := assert.New(t)

	t.Run("creating a transaction and getting all from a tracker", func(t *testing.T) {
		tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
		defer teardown()
		sharedTxns := NewSharedTxnsRepository(tbl)

		testInput := CreateSharedTxnInput{
			ID:           "test-shared-txn-id",
			TrackerID:    "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
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
	})

	t.Run("creating an unsettled transaction and a settled transaction and getting all unsettled transactions from a tracker", func(t *testing.T) {
		tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
		defer teardown()
		sharedTxns := NewSharedTxnsRepository(tbl)

		testInput := CreateSharedTxnInput{
			ID:           "test-shared-txn-id",
			TrackerID:    "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}
		err := sharedTxns.Create(testInput)
		assert.NoError(err)

		testUnsettledInput := CreateSharedTxnInput{
			ID:           "test-unsettled-shared-txn-id",
			TrackerID:    "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
			Unsettled:    true,
		}
		err = sharedTxns.Create(testUnsettledInput)
		assert.NoError(err)

		got, err := sharedTxns.GetUnsettledFromTracker("test-tracker-id")
		assert.NoError(err)

		want := SharedTxnItem{
			PK:           makeTrackerIDKey(testUnsettledInput.TrackerID),
			SK:           makeSharedTxnIDKey(testUnsettledInput.ID),
			EntityType:   sharedTxnEntityType,
			ID:           testUnsettledInput.ID,
			Tracker:      testUnsettledInput.TrackerID,
			Participants: testUnsettledInput.Participants,
			Unsettled:    unsettledFlagTrue,
		}
		assert.Contains(got, want)
	})
}
