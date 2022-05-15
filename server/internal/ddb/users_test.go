package ddb

import (
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/stretchr/testify/assert"
)

const testUsersTableName = "expenseus-testing-users"

func TestUsersTable(t *testing.T) {
	assert := assert.New(t)
	tbl, teardown := SetUpTestTable(t, testUsersTableName)
	defer teardown()
	users := NewUserRepository(tbl)

	// retrieving a non-existent user will give an error
	_, err := users.Get("non-existent-user")
	assert.EqualError(err, table.ErrItemNotFound.Error())

	user := UserItem{
		PK:         "user#test",
		SK:         "user#test",
		EntityType: "user",
		ID:         "test",
		GSI1PK:     "users",
		GSI1SK:     "user#test",
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
