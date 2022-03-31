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
}

type UsersTable interface {
	Get(id string) (UserItem, error)
	// GetAll() ([]UserItem, error)
	PutIfNotExists(item UserItem) error
	Delete(id string) error
}

type users struct {
	table *table.Table
}

const (
	HashKey       = "PK"
	RangeKey      = "SK"
	UserKeyPrefix = "user"
)

func NewUsersTable(t *table.Table) UsersTable {
	t.WithHashKey(HashKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(RangeKey, dynamodb.ScalarAttributeTypeS)
	return &users{table: t}
}

func (u *users) Get(id string) (UserItem, error) {
	userKey := fmt.Sprintf("%s#%s", UserKeyPrefix, id)
	item := &UserItem{}
	err := u.table.GetItem(attributes.String(userKey), attributes.String(userKey), item)
	if err != nil {
		return UserItem{}, err
	}
	return *item, nil
}

// func (u *users) GetAll() ([]UserItem, error) {
// 	response, err := u.table.DynamoDB.Scan(&dynamodb.ScanInput{TableName: u.table.Name})
// 	if err != nil {
// 		return nil, err
// 	}
// 	var items []UserItem

// 	for _, i := range response.Items {
// 		var item UserItem
// 		err = dynamodbattribute.UnmarshalMap(i, &item)
// 		if err != nil {
// 			return nil, err
// 		}
// 		items = append(items, item)
// 	}

// 	return items, nil
// }

func (u *users) PutIfNotExists(item UserItem) error {
	err := u.table.PutItem(item, option.PutCondition("attribute_not_exists(SK)"))
	if err != nil {
		return conflictOrErr(err)
	}

	return nil
}

func (u *users) Delete(id string) error {
	userKey := fmt.Sprintf("%s#%s", UserKeyPrefix, id)
	return u.table.DeleteItem(attributes.String(userKey), attributes.String(userKey))
}
