package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
)

const (
	trackerKeyPrefix  = "tracker"
	trackerEntityType = "tracker"
	allTrackersKey    = "trackers"
)

type TrackerItem struct {
	PK         string   `json:"PK"`
	SK         string   `json:"SK"`
	EntityType string   `json:"EntityType"`
	ID         string   `json:"ID"`
	Name       string   `json:"Name"`
	Users      []string `json:"Users"`
	GSI1PK     string   `json:"GSI1PK"`
	GSI1SK     string   `json:"GSI1SK"`
}

type CreateTrackerInput struct {
	ID    string
	Name  string
	Users []string
}

type TrackersRepository interface {
	Get(id string) (*TrackerItem, error)
	GetAll() ([]TrackerItem, error)
	Put(item TrackerItem) error
	Create(input CreateTrackerInput) error
}

type trackersRepo struct {
	table *table.Table
}

func NewTrackersRepository(t *table.Table) TrackersRepository {
	t.WithHashKey(tablePrimaryKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(tableSortKey, dynamodb.ScalarAttributeTypeS)
	return &trackersRepo{t}
}

func (t *trackersRepo) Get(id string) (*TrackerItem, error) {
	// make the id
	trackerIDKey := fmt.Sprintf("%s#%s", trackerKeyPrefix, id)
	item := &TrackerItem{}
	err := t.table.GetItem(attributes.String(trackerIDKey), attributes.String(trackerIDKey), item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *trackersRepo) Put(item TrackerItem) error {
	return t.table.PutItem(item)
}

func (t *trackersRepo) Create(input CreateTrackerInput) error {
	trackerIDKey := fmt.Sprintf("%s#%s", trackerKeyPrefix, input.ID)

	return t.table.PutItem(TrackerItem{
		PK:         trackerIDKey,
		SK:         trackerIDKey,
		EntityType: trackerEntityType,
		ID:         input.ID,
		Name:       input.Name,
		Users:      input.Users,
		GSI1PK:     allTrackersKey,
		GSI1SK:     trackerIDKey,
	})
}

func (t *trackersRepo) GetAll() ([]TrackerItem, error) {
	options := []option.QueryInput{
		option.Index("GSI1"),
		option.QueryExpressionAttributeName(gsi1PrimaryKey, "#GSI1PK"),
		option.QueryExpressionAttributeValue(":allTrackersKey", attributes.String(allTrackersKey)),
		option.QueryKeyConditionExpression("#GSI1PK = :allTrackersKey"),
	}

	var items []TrackerItem

	_, err := t.table.Query(&items, options...)
	if err != nil {
		return nil, err
	}

	return items, nil
}
