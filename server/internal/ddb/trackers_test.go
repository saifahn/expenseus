package ddb

import (
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/stretchr/testify/assert"
)

const testTrackersTableName = "expenseus-testing-trackers"

func TestTrackersRepo(t *testing.T) {
	assert := assert.New(t)
	dynamoDB := NewDynamoDBLocalAPI()

	err := CreateTestTable(dynamoDB, testTrackersTableName)
	if err != nil {
		t.Logf("table could not be created: %v", err)
	}
	defer DeleteTable(dynamoDB, testTrackersTableName)
	tbl := table.New(dynamoDB, testTrackersTableName)
	trackers := NewTrackersRepository(tbl)

	_, err = trackers.Get("non-existent-item")
	assert.EqualError(err, table.ErrItemNotFound.Error())

	item := &TrackerItem{
		PK:         "tracker#test-tracker-id",
		SK:         "tracker#test-tracker-id",
		EntityType: trackerEntityType,
		ID:         "test-tracker-id",
	}

	err = trackers.Put(*item)
	assert.NoError(err)

	got, err := trackers.Get(item.ID)
	assert.NoError(err)
	assert.Equal(item, got)
}
