package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
)

const (
	trackerKeyPrefix  = "tracker"
	trackerEntityType = "tracker"
)

type TrackerItem struct {
	PK         string   `json:"PK"`
	SK         string   `json:"SK"`
	EntityType string   `json:"EntityType"`
	ID         string   `json:"ID"`
	Name       string   `json:"Name"`
	Users      []string `json:"Users"`
}

type CreateTrackerInput struct {
	ID    string
	Name  string
	Users []string
}

type TrackersRepository interface {
	Get(id string) (*TrackerItem, error)
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
	})
}
