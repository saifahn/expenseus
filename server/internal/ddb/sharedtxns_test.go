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
			ID:           id,
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}

		err := sharedTxns.Create(testInput)
		assert.NoError(err)

		got, err := sharedTxns.GetFromTracker(testInput.Tracker)
		assert.NoError(err)

		want := SharedTxnItem{
			PK:           makeTrackerIDKey(testInput.Tracker),
			SK:           makeSharedTxnIDKey(id),
			GSI1PK:       makeTrackerIDKey(testInput.Tracker),
			GSI1SK:       makeSharedTxnDateIDKey(testInput),
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

		testInput := app.SharedTransaction{
			ID:           "test-shared-txn-id",
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}
		err := sharedTxns.Create(testInput)
		assert.NoError(err)

		unsettledID := "test-unsettled-txn-id"
		testUnsettledInput := app.SharedTransaction{
			ID:           unsettledID,
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
			Unsettled:    true,
		}
		err = sharedTxns.Create(testUnsettledInput)
		assert.NoError(err)

		got, err := sharedTxns.GetUnsettledFromTracker("test-tracker-id")
		assert.NoError(err)

		want := SharedTxnItem{
			PK:           makeTrackerIDKey(testUnsettledInput.Tracker),
			SK:           makeSharedTxnIDKey(unsettledID),
			GSI1PK:       makeTrackerIDKey(testUnsettledInput.Tracker),
			GSI1SK:       makeSharedTxnDateIDKey(testUnsettledInput),
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

		testUnsettledInput := app.SharedTransaction{
			ID:           "test-unsettled-shared-txn-id",
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
			Unsettled:    true,
		}
		err := sharedTxns.Create(testUnsettledInput)
		assert.NoError(err)

		testSettleTxnPayload := SettleTxnInput{
			ID:           "test-unsettled-shared-txn-id",
			TrackerID:    "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}

		err = sharedTxns.Settle([]SettleTxnInput{testSettleTxnPayload})
		assert.NoError(err)

		got, err := sharedTxns.GetUnsettledFromTracker("test-tracker-id")
		assert.NoError(err)
		assert.Empty(got)
	})

	t.Run("updating an transaction", func(t *testing.T) {
		tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
		defer teardown()
		sharedTxns := NewSharedTxnsRepository(tbl)

		initialTxn := app.SharedTransaction{
			ID:           "test-shared-txn-id",
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
			Amount:       123,
		}
		err := sharedTxns.Create(initialTxn)
		assert.NoError(err)

		updatedTxn := app.SharedTransaction{
			ID:           "test-shared-txn-id",
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
			Amount:       456,
		}
		err = sharedTxns.Update(updatedTxn)
		assert.NoError(err)

		got, err := sharedTxns.GetFromTracker(initialTxn.Tracker)
		assert.NoError(err)
		want := SharedTxnItem{
			PK:           makeTrackerIDKey(updatedTxn.Tracker),
			SK:           makeSharedTxnIDKey(updatedTxn.ID),
			GSI1PK:       makeTrackerIDKey(updatedTxn.Tracker),
			GSI1SK:       makeSharedTxnDateIDKey(updatedTxn),
			EntityType:   sharedTxnEntityType,
			ID:           updatedTxn.ID,
			Tracker:      updatedTxn.Tracker,
			Participants: updatedTxn.Participants,
			Amount:       updatedTxn.Amount,
		}
		assert.ElementsMatch(got, []SharedTxnItem{want})
	})

	t.Run("deleting a transaction", func(t *testing.T) {
		tbl, teardown := SetUpTestTable(t, testSharedTxnsTableName)
		defer teardown()
		sharedTxns := NewSharedTxnsRepository(tbl)

		id := "test-shared-txn-id"
		testInput := app.SharedTransaction{
			ID:           id,
			Tracker:      "test-tracker-id",
			Participants: []string{"test-01", "test-02"},
		}
		err := sharedTxns.Create(testInput)
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

func TestGetSharedBetweenDates(t *testing.T) {
	testTrackerID := "test-tracker-id"
	testTxnID := "test-txn-id"

	initialTxn := app.SharedTransaction{
		ID:           testTxnID,
		Tracker:      testTrackerID,
		Date:         10000,
		Participants: []string{"user-01", "user-02"},
	}

	wantItem := SharedTxnItem{
		PK:           makeTrackerIDKey(testTrackerID),
		SK:           makeSharedTxnIDKey(testTxnID),
		GSI1PK:       makeTrackerIDKey(testTrackerID),
		GSI1SK:       "txn.shared#10000#test-txn-id",
		Participants: []string{"user-01", "user-02"},
		EntityType:   sharedTxnEntityType,
		ID:           testTxnID,
		Tracker:      testTrackerID,
		Date:         10000,
	}

	tests := map[string]struct {
		wantItems []SharedTxnItem
		from      int64
		to        int64
	}{
		"with a date-range containing an item": {
			wantItems: []SharedTxnItem{wantItem},
			from:      10000,
			to:        20000,
		},
		"with a date range not containing any items": {
			wantItems: []SharedTxnItem{},
			from:      15000,
			to:        20000,
		},
	}
	for name, tc := range tests {
		assert := assert.New(t)

		t.Run(name, func(t *testing.T) {
			tbl, teardown := SetUpTestTable(t, "test-get-shared-txns-between-dates")
			defer teardown()
			shared := NewSharedTxnsRepository(tbl)

			err := shared.Create(initialTxn)
			assert.NoError(err)

			got, err := shared.GetFromTrackerBetweenDates(wantItem.Tracker, tc.from, tc.to)
			assert.NoError(err)
			assert.ElementsMatch(got, tc.wantItems)
		})
	}
}
