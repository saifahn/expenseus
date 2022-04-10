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
}
