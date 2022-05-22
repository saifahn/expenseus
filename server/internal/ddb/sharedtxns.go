package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/table/option"
	"github.com/saifahn/expenseus/internal/app"
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
	Date         int64    `json:"Date"`
	Amount       int64    `json:"Amount"`
	Shop         string   `json:"Shop"`
	Tracker      string   `json:"Tracker"`
	Participants []string `json:"Participants"`
	Unsettled    string   `json:"Unsettled,omitempty"`
}
type SettleTxnInputItem struct {
	ID           string
	TrackerID    string
	Participants []string
}

type SharedTxnsRepository interface {
	Create(txnID string, input app.SharedTransaction) error
	GetFromTracker(trackerID string) ([]SharedTxnItem, error)
	GetUnsettledFromTracker(trackerID string) ([]SharedTxnItem, error)
	Settle(input []SettleTxnInputItem) error
}

type sharedTxnsRepo struct {
	table *table.Table
}

func NewSharedTxnsRepository(t *table.Table) SharedTxnsRepository {
	t.WithHashKey(tablePrimaryKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(tableSortKey, dynamodb.ScalarAttributeTypeS)
	return &sharedTxnsRepo{t}
}

func (r *sharedTxnsRepo) Create(txnID string, input app.SharedTransaction) error {
	trackerIDKey := makeTrackerIDKey(input.Tracker)
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
			ID:           txnID,
			Tracker:      input.Tracker,
			Participants: input.Participants,
			Unsettled:    unsettledVal,
			Date:         input.Date,
			Amount:       input.Amount,
			Shop:         input.Shop,
		})
		if err != nil {
			return err
		}
	}

	err := r.table.PutItem(SharedTxnItem{
		PK:           trackerIDKey,
		SK:           txnIDKey,
		EntityType:   sharedTxnEntityType,
		ID:           txnID,
		Tracker:      input.Tracker,
		Participants: input.Participants,
		Unsettled:    unsettledVal,
		Date:         input.Date,
		Amount:       input.Amount,
		Shop:         input.Shop,
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

// Settle takes a slice of SettleTxnPayload and removes the "Unsettled"
// attribute from the database items that correspond to the transactions
func (r *sharedTxnsRepo) Settle(txns []SettleTxnInputItem) error {
	var updateOpts []option.UpdateItemInput
	updateOpts = append(updateOpts, option.UpdateExpressionAttributeName("Unsettled", "#unsettled"))
	updateOpts = append(updateOpts, option.UpdateExpression("REMOVE #unsettled"))

	for _, t := range txns {
		trackerIDKey := makeTrackerIDKey(t.TrackerID)
		txnIDKey := makeSharedTxnIDKey(t.ID)
		_, err := r.table.UpdateItem(
			attributes.String(trackerIDKey),
			attributes.String(txnIDKey),
			updateOpts...,
		)
		if err != nil {
			return err
		}

		for _, u := range t.Participants {
			userKey := makeUserIDKey(u)
			_, err := r.table.UpdateItem(
				attributes.String(userKey),
				attributes.String(txnIDKey),
				updateOpts...,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func makeSharedTxnIDKey(id string) string {
	return fmt.Sprintf("%s#%s", sharedTxnKeyPrefix, id)
}
