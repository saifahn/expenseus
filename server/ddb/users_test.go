package ddb

import (
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/saifahn/expenseus"
	"github.com/stretchr/testify/assert"
)

const testUsersTableName = "expenseus-testing-users"

func TestUsersTable(t *testing.T) {
	assert := assert.New(t)
	dynamodb := NewDynamoDBLocalAPI()

	err := CreateTestTable(dynamodb, testUsersTableName)
	if err != nil {
		t.Logf("table could not be created: %v", err)
	}
	defer DeleteTable(dynamodb, testUsersTableName)

	tbl := table.New(dynamodb, testUsersTableName)
	users := NewUsersTable(tbl)

	// retrieving a non-existent user will give an error
	_, err = users.Get("non-existent-user")
	assert.EqualError(err, table.ErrItemNotFound.Error())

	user := UserItem{
		User: expenseus.User{
			ID:       "test-user",
			Name:     "Testman",
			Username: "testman-23",
		},
	}

	err = users.PutIfNotExists(user)
	assert.NoError(err)

	// trying to put the same user will result in an error
	err = users.PutIfNotExists(user)
	assert.EqualError(err, ErrConflict.Error())

	// the user can be retrieved
	got, err := users.Get(user.ID)
	assert.NoError(err)
	assert.Equal(user, got)

	// retrieve all users
	usersGot, err := users.GetAll()
	assert.NoError(err)
	assert.Len(usersGot, 1)
	assert.Contains(usersGot, user)

	err = users.Delete(user.ID)
	assert.NoError(err)
	_, err = users.Get(user.ID)
	assert.EqualError(err, table.ErrItemNotFound.Error())
}
