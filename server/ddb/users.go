package ddb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
	"github.com/saifahn/expenseus"
)

type UserItem struct {
	expenseus.User
}

type UsersTable interface {
	Get(id string) (UserItem, error)
	GetAll() ([]UserItem, error)
	PutIfNotExists(item UserItem) error
	Delete(id string) error
}

type usersTable struct {
	table *table.Table
}

const UsersHashKeyName = "id"

func NewUsersTable(t *table.Table) UsersTable {
	t.WithHashKey(UsersHashKeyName, dynamodb.ScalarAttributeTypeS)
	return &usersTable{table: t}
}

func (u *usersTable) Get(id string) (UserItem, error) {
	item := &UserItem{}
	err := u.table.GetItem(attributes.String(id), nil, item, option.ConsistentRead())
	if err != nil {
		return UserItem{}, err
	}
	return *item, nil
}

func (u *usersTable) GetAll() ([]UserItem, error) {
	response, err := u.table.DynamoDB.Scan(&dynamodb.ScanInput{TableName: u.table.Name})
	if err != nil {
		return nil, err
	}
	var items []UserItem

	for _, i := range response.Items {
		var item UserItem
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (u *usersTable) PutIfNotExists(item UserItem) error {
	err := u.table.PutItem(item, option.PutCondition("attribute_not_exists(id)"))
	if err != nil {
		return conflictOrErr(err)
	}

	return nil
}

func (u *usersTable) Delete(id string) error {
	return u.table.DeleteItem(attributes.String(id), nil)
}
