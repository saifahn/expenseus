package ddb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testTrackersTableName = "expenseus-testing-trackers"

func TestCreateTracker(t *testing.T) {
	assert := assert.New(t)
	tbl, teardown := SetUpTestTable(t, testTrackersTableName)
	defer teardown()
	trackers := NewTrackersRepository(tbl)

	testUserIDs := []string{"test-01", "test-02"}
	testInput := CreateTrackerInput{
		Users: testUserIDs,
		ID:    "test-tracker-id",
		Name:  "The Test Tracker",
	}

	err := trackers.Create(testInput)
	assert.NoError(err)

	got, err := trackers.Get(testInput.ID)
	assert.NoError(err)

	want := &TrackerItem{
		PK:         "tracker#test-tracker-id",
		SK:         "tracker#test-tracker-id",
		EntityType: trackerEntityType,
		ID:         "test-tracker-id",
		Name:       "The Test Tracker",
		Users:      testUserIDs,
		GSI1PK:     allTrackersKey,
		GSI1SK:     "tracker#test-tracker-id",
	}
	assert.Equal(want, got)

	// allGot, err := trackers.GetAll()
	// assert.NoError(err)

	// allExpected := []TrackerItem{*want}
	// assert.ElementsMatch(allExpected, allGot)
}
