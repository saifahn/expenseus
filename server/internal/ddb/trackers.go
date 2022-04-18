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
	Create(input CreateTrackerInput) error
	Get(id string) (*TrackerItem, error)
	GetAll() ([]TrackerItem, error)
	GetByUser(userID string) ([]TrackerItem, error)
}

type trackersRepo struct {
	table *table.Table
}

func NewTrackersRepository(t *table.Table) TrackersRepository {
	t.WithHashKey(tablePrimaryKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(tableSortKey, dynamodb.ScalarAttributeTypeS)
	return &trackersRepo{t}
}

func (t *trackersRepo) Create(input CreateTrackerInput) error {
	trackerIDKey := makeTrackerIDKey(input.ID)
	userIDKey := makeUserIDKey(input.Users[0])

	err := t.table.PutItem(TrackerItem{
		PK:         trackerIDKey,
		SK:         trackerIDKey,
		EntityType: trackerEntityType,
		ID:         input.ID,
		Name:       input.Name,
		Users:      input.Users,
		GSI1PK:     allTrackersKey,
		GSI1SK:     trackerIDKey,
	})
	if err != nil {
		return err
	}

	err = t.table.PutItem(TrackerItem{
		PK:         userIDKey,
		SK:         trackerIDKey,
		EntityType: trackerEntityType,
		ID:         input.ID,
		Name:       input.Name,
		Users:      input.Users,
		GSI1PK:     allTrackersKey,
		GSI1SK:     trackerIDKey,
	})
	return err
}

func (t *trackersRepo) Get(id string) (*TrackerItem, error) {
	trackerIDKey := makeTrackerIDKey(id)

	item := &TrackerItem{}
	err := t.table.GetItem(attributes.String(trackerIDKey), attributes.String(trackerIDKey), item)
	if err != nil {
		return nil, err
	}
	return item, nil
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

func (t *trackersRepo) GetByUser(userID string) ([]TrackerItem, error) {
	userIDKey := makeUserIDKey(userID)
	allTrackerPrefix := fmt.Sprintf("%s#", trackerKeyPrefix)

	options := []option.QueryInput{
		option.QueryExpressionAttributeName(tablePrimaryKey, "#PK"),
		option.QueryExpressionAttributeName(tableSortKey, "#SK"),
		option.QueryExpressionAttributeValue(":userKey", attributes.String(userIDKey)),
		option.QueryExpressionAttributeValue(":allTrackerPrefix", attributes.String(allTrackerPrefix)),
		option.QueryKeyConditionExpression("#PK = :userKey and begins_with(#SK, :allTrackerPrefix)"),
	}

	var items []TrackerItem
	_, err := t.table.Query(&items, options...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func makeTrackerIDKey(id string) string {
	return fmt.Sprintf("%s#%s", trackerKeyPrefix, id)
}
