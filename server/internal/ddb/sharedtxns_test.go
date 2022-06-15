package ddb

import (
	"testing"

	"github.com/saifahn/expenseus/internal/app"
	"github.com/stretchr/testify/assert"
)

const testSharedTxnsTableName = "expenseus-testing-shared-txns"

func TestSharedTxns(t *testing.T) {
	assert := assert.New(t)

	t.Run("creating a transaction and getting all from a tracker", func(t *testing.T) {
		tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
		defer teardown()
		sharedTxns := NewSharedTxnsRepository(tbl)

		id := "test-shared-txn-id"
		testInput := app.SharedTransaction{
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}

		err := sharedTxns.Create(id, testInput)
		assert.NoError(err)

		got, err := sharedTxns.GetFromTracker(testInput.Tracker)
		assert.NoError(err)

		want := SharedTxnItem{
			PK:           makeTrackerIDKey(testInput.Tracker),
			SK:           makeSharedTxnIDKey(id),
			EntityType:   sharedTxnEntityType,
			ID:           id,
			Tracker:      testInput.Tracker,
			Participants: testInput.Participants,
		}
		assert.Contains(got, want)
	})

	t.Run("creating an unsettled transaction and a settled transaction and getting all unsettled transactions from a tracker", func(t *testing.T) {
		tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
		defer teardown()
		sharedTxns := NewSharedTxnsRepository(tbl)

		id := "test-shared-txn-id"
		testInput := app.SharedTransaction{
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}
		err := sharedTxns.Create(id, testInput)
		assert.NoError(err)

		unsettledID := "test-unsettled-shared-txn-id"
		testUnsettledInput := app.SharedTransaction{
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
			Unsettled:    true,
		}
		err = sharedTxns.Create(unsettledID, testUnsettledInput)
		assert.NoError(err)

		got, err := sharedTxns.GetUnsettledFromTracker("test-tracker-id")
		assert.NoError(err)

		want := SharedTxnItem{
			PK:           makeTrackerIDKey(testUnsettledInput.Tracker),
			SK:           makeSharedTxnIDKey(unsettledID),
			EntityType:   sharedTxnEntityType,
			ID:           unsettledID,
			Tracker:      testUnsettledInput.Tracker,
			Participants: testUnsettledInput.Participants,
			Unsettled:    unsettledFlagTrue,
		}
		assert.Contains(got, want)
	})

	t.Run("creating an unsettled transaction and marking all transactions in a tracker settled", func(t *testing.T) {
		tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
		defer teardown()
		sharedTxns := NewSharedTxnsRepository(tbl)

		unsettledID := "test-unsettled-shared-txn-id"
		testUnsettledInput := app.SharedTransaction{
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
			Unsettled:    true,
		}
		err := sharedTxns.Create(unsettledID, testUnsettledInput)
		assert.NoError(err)

		testSettleTxnPayload := SettleTxnInput{
			ID:           unsettledID,
			TrackerID:    "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}

		err = sharedTxns.Settle([]SettleTxnInput{testSettleTxnPayload})
		assert.NoError(err)

		got, err := sharedTxns.GetUnsettledFromTracker("test-tracker-id")
		assert.NoError(err)
		assert.Empty(got)
	})

	t.Run("deleting a transaction", func(t *testing.T) {
		tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
		defer teardown()
		sharedTxns := NewSharedTxnsRepository(tbl)

		id := "test-shared-txn-id"
		testInput := app.SharedTransaction{
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}
		err := sharedTxns.Create(id, testInput)
		assert.NoError(err)

		testDelInput := app.DelSharedTxnInput{
			TxnID:        id,
			Tracker:      testInput.Tracker,
			Participants: testInput.Participants,
		}
		err = sharedTxns.Delete(testDelInput)
		assert.NoError(err)

		got, err := sharedTxns.GetFromTracker(testInput.Tracker)
		assert.NoError(err)
		assert.Empty(got)
	})
}
