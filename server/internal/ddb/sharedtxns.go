package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/table"
)

const (
	sharedTxnKeyPrefix  = "txn.shared"
	sharedTxnEntityType = "sharedTransaction"
)

type SharedTxnItem struct {
	PK           string   `json:"PK"`
	SK           string   `json:"SK"`
	EntityType   string   `json:"EntityType"`
	ID           string   `json:"ID"`
	Tracker      string   `json:"Tracker"`
	Participants []string `json:"Participants"`
}

type CreateSharedTxnInput struct {
	ID           string
	TrackerID    string
	Participants []string
}

type SharedTxnsRepository interface {
	Create(input CreateSharedTxnInput) error
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
	for _, p := range input.Participants {
		userIDKey := makeUserIDKey(p)
		err := r.table.PutItem(SharedTxnItem{
			PK:           userIDKey,
			SK:           txnIDKey,
			EntityType:   sharedTxnEntityType,
			ID:           input.ID,
			Tracker:      input.TrackerID,
			Participants: input.Participants,
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
	})
	return err
}

func makeSharedTxnIDKey(id string) string {
	return fmt.Sprintf("%s#%s", sharedTxnKeyPrefix, id)
}
