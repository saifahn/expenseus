package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
)

const (
	sharedTxnKeyPrefix  = "txn.shared"
	sharedTxnEntityType = "sharedTransaction"
	unsettledFlagTrue   = "X"
)

type SharedTxnItem struct {
	PK           string   `json:"PK"`
	SK           string   `json:"SK"`
	EntityType   string   `json:"EntityType"`
	ID           string   `json:"ID"`
	Tracker      string   `json:"Tracker"`
	Participants []string `json:"Participants"`
	Unsettled    string   `json:"Unsettled,omitempty"`
}

type CreateSharedTxnInput struct {
	ID           string
	TrackerID    string
	Participants []string
	Unsettled    bool
}

type SharedTxnsRepository interface {
	Create(input CreateSharedTxnInput) error
	GetFromTracker(trackerID string) ([]SharedTxnItem, error)
	GetUnsettledFromTracker(trackerID string) ([]SharedTxnItem, error)
}

type sharedTxnsRepo struct {
	table *table.Table
}

func NewSharedTxnsRepository(t *table.Table) SharedTxnsRepository {
	t.WithHashKey(tablePrimaryKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(tableSortKey, dynamodb.ScalarAttributeTypeS)
	return &sharedTxnsRepo{t}
}

func (r *sharedTxnsRepo) Create(input CreateSharedTxnInput) error {
	trackerIDKey := makeTrackerIDKey(input.TrackerID)
	txnIDKey := makeSharedTxnIDKey(input.ID)
	var unsettledVal string
	if input.Unsettled {
		unsettledVal = unsettledFlagTrue
	}

	for _, p := range input.Participants {
		userIDKey := makeUserIDKey(p)
		err := r.table.PutItem(SharedTxnItem{
			PK:           userIDKey,
			SK:           txnIDKey,
			EntityType:   sharedTxnEntityType,
			ID:           input.ID,
			Tracker:      input.TrackerID,
			Participants: input.Participants,
			Unsettled:    unsettledVal,
		})
		if err != nil {
			return err
		}
	}

	err := r.table.PutItem(SharedTxnItem{
		PK:           trackerIDKey,
		SK:           txnIDKey,
		EntityType:   sharedTxnEntityType,
		ID:           input.ID,
		Tracker:      input.TrackerID,
		Participants: input.Participants,
		Unsettled:    unsettledVal,
	})
	return err
}

func (r *sharedTxnsRepo) GetFromTracker(trackerID string) ([]SharedTxnItem, error) {
	trackerIDKey := makeTrackerIDKey(trackerID)

	options := []option.QueryInput{
		option.QueryExpressionAttributeName(tablePrimaryKey, "#PK"),
		option.QueryExpressionAttributeName(tableSortKey, "#SK"),
		option.QueryExpressionAttributeValue(":trackerID", attributes.String(trackerIDKey)),
		option.QueryExpressionAttributeValue(":sharedTxnKeyPrefix", attributes.String(sharedTxnKeyPrefix)),
		option.QueryKeyConditionExpression("#PK = :trackerID and begins_with(#SK, :sharedTxnKeyPrefix)"),
	}

	var items []SharedTxnItem

	_, err := r.table.Query(&items, options...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *sharedTxnsRepo) GetUnsettledFromTracker(trackerID string) ([]SharedTxnItem, error) {
	trackerIDKey := makeTrackerIDKey(trackerID)

	options := []option.QueryInput{
		option.Index(unsettledTxnsIndexName),
		option.QueryExpressionAttributeName(unsettledTxnsIndexPK, "#unsettledPK"),
		option.QueryExpressionAttributeName(unsettledTxnsIndexSK, "#unsettledSK"),
		option.QueryExpressionAttributeValue(":trackerID", attributes.String(trackerIDKey)),
		option.QueryExpressionAttributeValue(":true", attributes.String(unsettledFlagTrue)),
		option.QueryKeyConditionExpression("#unsettledPK = :trackerID and #unsettledSK = :true"),
	}

	var items []SharedTxnItem

	_, err := r.table.Query(&items, options...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func makeSharedTxnIDKey(id string) string {
	return fmt.Sprintf("%s#%s", sharedTxnKeyPrefix, id)
}
