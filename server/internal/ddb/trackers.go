package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
)

const (
	trackerKeyPrefix = "tracker"
)

type TrackerItem struct {
	PK string `json:"PK"`
	SK string `json:"SK"`
}

type TrackersRepository interface {
	Get(id string) (TrackerItem, error)
}

type trackersRepo struct {
	table *table.Table
}

func NewTrackersRepository(t *table.Table) TrackersRepository {
	t.WithHashKey(tablePrimaryKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(tableSortKey, dynamodb.ScalarAttributeTypeS)
	return &trackersRepo{t}
}

func (t *trackersRepo) Get(id string) (TrackerItem, error) {
	// make the id
	trackerIDKey := fmt.Sprintf("%s#%s", trackerKeyPrefix, id)
	item := &TrackerItem{}
	err := t.table.GetItem(attributes.String(trackerIDKey), attributes.String(trackerIDKey), item)
	if err != nil {
		return TrackerItem{}, err
	}
	return *item, nil
}
