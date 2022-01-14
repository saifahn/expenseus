package dynamodb

import (
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/stretchr/testify/assert"
)

const testUsersTableName = "expenseus-testing-users"

func TestUsersTable(t *testing.T) {
	assert := assert.New(t)
	dynamodb := newDynamoDBLocalAPI()

	err := createTestTable(dynamodb, testUsersTableName)
	if err != nil {
		t.Logf("table could not be created: %v", err)
	}
	tbl := table.New(dynamodb, testUsersTableName)
	users := NewUsersTable(tbl)

	// retrieving a non-existent user will give an error
	_, err = users.Get("non-existent-user")
	assert.EqualError(err, table.ErrItemNotFound.Error())

	user := &UserItem{
		ID:           "test-user",
		EmailAddress: "test-user@test.com",
		ExternalID:   "test-external-id",
	}

	err = users.PutIfNotExists(*user)
	assert.NoError(err)
}
