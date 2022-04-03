package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
)

type UserItem struct {
	PK         string `json:"PK"`
	SK         string `json:"SK"`
	EntityType string `json:"EntityType"`
	ID         string `json:"ID"`
	GSI1PK     string `json:"GSI1PK"`
	GSI1SK     string `json:"GSI1SK"`
}

type UsersTable interface {
	Get(id string) (UserItem, error)
	GetAll() ([]UserItem, error)
	PutIfNotExists(item UserItem) error
	Delete(id string) error
}

type users struct {
	table *table.Table
}

const (
	userKeyPrefix  = "user"
	userEntityType = "user"
	allUsersKey    = "users"
)

func NewUsersTable(t *table.Table) UsersTable {
	t.WithHashKey(tablePrimaryKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(tableSortKey, dynamodb.ScalarAttributeTypeS)
	return &users{table: t}
}

func (u *users) Get(id string) (UserItem, error) {
	key := makeUserIDKey(id)
	item := &UserItem{}
	err := u.table.GetItem(attributes.String(key), attributes.String(key), item)
	if err != nil {
		return UserItem{}, err
	}
	return *item, nil
}

func (u *users) GetAll() ([]UserItem, error) {
	options := []option.QueryInput{
		option.Index("GSI1"),
		option.QueryExpressionAttributeName(gsi1PrimaryKey, "#GSI1PK"),
		option.QueryExpressionAttributeValue(":usersKey", attributes.String(allUsersKey)),
		option.QueryKeyConditionExpression("#GSI1PK = :usersKey"),
	}

	var items []UserItem

	_, err := u.table.Query(&items, options...)

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (u *users) PutIfNotExists(item UserItem) error {
	err := u.table.PutItem(item, option.PutCondition("attribute_not_exists(SK)"))
	if err != nil {
		return conflictOrErr(err)
	}

	return nil
}

func (u *users) Delete(id string) error {
	key := makeUserIDKey(id)
	return u.table.DeleteItem(attributes.String(key), attributes.String(key))
}

func makeUserIDKey(id string) string {
	return fmt.Sprintf("%s#%s", userKeyPrefix, id)
}
