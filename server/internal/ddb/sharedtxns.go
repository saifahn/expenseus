package ddb

import (
	"encoding/json"
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
	GSI1PK       string   `json:"GSI1PK"`
	GSI1SK       string   `json:"GSI1SK"`
	EntityType   string   `json:"EntityType"`
	ID           string   `json:"ID"`
	Date         int64    `json:"Date"`
	Amount       int64    `json:"Amount"`
	Location     string   `json:"Location"`
	Tracker      string   `json:"Tracker"`
	Category     string   `json:"Category"`
	Participants []string `json:"Participants"`
	Payer        string   `json:"Payer"`
	Unsettled    string   `json:"Unsettled,omitempty"`
	Details      string   `json:"Details"`
	SplitJSON    string   `json:"Split"`
}

type SettleTxnInput struct {
	ID           string
	TrackerID    string
	Participants []string
}

type SharedTxnsRepository interface {
	Create(txn app.SharedTransaction) error
	GetFromTracker(trackerID string) ([]SharedTxnItem, error)
	GetFromTrackerBetweenDates(trackerID string, from, to int64) ([]SharedTxnItem, error)
	GetByUserBetweenDates(userID string, from, to int64) ([]SharedTxnItem, error)
	GetUnsettledFromTracker(trackerID string) ([]SharedTxnItem, error)
	Update(txn app.SharedTransaction) error
	Delete(input app.DelSharedTxnInput) error
	Settle(input []SettleTxnInput) error
}

type sharedTxnsRepo struct {
	table *table.Table
}

func NewSharedTxnsRepository(t *table.Table) SharedTxnsRepository {
	t.WithHashKey(tablePrimaryKey, dynamodb.ScalarAttributeTypeS)
	t.WithRangeKey(tableSortKey, dynamodb.ScalarAttributeTypeS)
	return &sharedTxnsRepo{t}
}

func (r *sharedTxnsRepo) Create(txn app.SharedTransaction) error {
	trackerIDKey := makeTrackerIDKey(txn.Tracker)
	txnIDKey := makeSharedTxnIDKey(txn.ID)
	txnDateKey := makeSharedTxnDateIDKey(txn)

	txnItem := SharedTxnItem{
		SK:           txnIDKey,
		GSI1SK:       txnDateKey,
		EntityType:   sharedTxnEntityType,
		ID:           txn.ID,
		Category:     txn.Category,
		Tracker:      txn.Tracker,
		Participants: txn.Participants,
		Date:         txn.Date,
		Amount:       txn.Amount,
		Location:     txn.Location,
		Payer:        txn.Payer,
		Details:      txn.Details,
	}

	if txn.Unsettled {
		txnItem.Unsettled = unsettledFlagTrue
	}
	// store the split map as a JSON string
	if txn.Split != nil {
		splitJSON, err := json.Marshal(txn.Split)
		if err != nil {
			return err
		}
		txnItem.SplitJSON = string(splitJSON)
	}

	for _, p := range txn.Participants {
		userIDKey := makeUserIDKey(p)
		txnItem.PK = userIDKey
		txnItem.GSI1PK = userIDKey
		err := r.table.PutItem(txnItem)
		if err != nil {
			return err
		}
	}

	txnItem.GSI1PK = trackerIDKey
	txnItem.PK = trackerIDKey
	err := r.table.PutItem(txnItem)
	return err
}

func (r *sharedTxnsRepo) GetFromTracker(trackerID string) ([]SharedTxnItem, error) {
	trackerIDKey := makeTrackerIDKey(trackerID)

	options := []option.QueryInput{
		option.Index(gsi1Name),
		option.QueryExpressionAttributeName(gsi1PrimaryKey, "#GSI1PK"),
		option.QueryExpressionAttributeName(gsi1SortKey, "#GSI1SK"),
		option.QueryExpressionAttributeValue(":trackerIDKey", attributes.String(trackerIDKey)),
		option.QueryExpressionAttributeValue(":sharedTxnKeyPrefix", attributes.String(sharedTxnKeyPrefix)),
		// sort descending
		option.Reverse(),
		option.QueryKeyConditionExpression("#GSI1PK = :trackerIDKey and begins_with(#GSI1SK, :sharedTxnKeyPrefix)"),
	}

	var items []SharedTxnItem

	_, err := r.table.Query(&items, options...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *sharedTxnsRepo) GetFromTrackerBetweenDates(trackerID string, from, to int64) ([]SharedTxnItem, error) {
	trackerIDKey := makeTrackerIDKey(trackerID)
	txnDateFromKey := makeSharedTxnDateKey(from)
	txnDateToKey := makeSharedTxnDateKey(to)

	options := []option.QueryInput{
		option.Index(gsi1Name),
		option.QueryExpressionAttributeName(gsi1PrimaryKey, "#GSI1PK"),
		option.QueryExpressionAttributeName(gsi1SortKey, "#GSI1SK"),
		option.QueryExpressionAttributeValue(":trackerKey", attributes.String(trackerIDKey)),
		option.QueryExpressionAttributeValue(":txnDateFromKey", attributes.String(txnDateFromKey)),
		option.QueryExpressionAttributeValue(":txnDateToKey", attributes.String(txnDateToKey)),
		option.QueryKeyConditionExpression("#GSI1PK = :trackerKey and #GSI1SK BETWEEN :txnDateFromKey AND :txnDateToKey"),
	}

	var items []SharedTxnItem
	_, err := r.table.Query(&items, options...)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *sharedTxnsRepo) GetByUserBetweenDates(userID string, from, to int64) ([]SharedTxnItem, error) {
	userIDKey := makeUserIDKey(userID)
	txnDateFromKey := makeSharedTxnDateKey(from)
	txnDateToKey := makeSharedTxnDateKey(to)

	options := []option.QueryInput{
		option.Index(gsi1Name),
		option.QueryExpressionAttributeName(gsi1PrimaryKey, "#GSI1PK"),
		option.QueryExpressionAttributeName(gsi1SortKey, "#GSI1SK"),
		option.QueryExpressionAttributeValue(":userKey", attributes.String(userIDKey)),
		option.QueryExpressionAttributeValue(":txnDateFromKey", attributes.String(txnDateFromKey)),
		option.QueryExpressionAttributeValue(":txnDateToKey", attributes.String(txnDateToKey)),
		option.QueryKeyConditionExpression("#GSI1PK = :userKey and #GSI1SK BETWEEN :txnDateFromKey AND :txnDateToKey"),
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

func (r *sharedTxnsRepo) Update(txn app.SharedTransaction) error {
	trackerIDKey := makeTrackerIDKey(txn.Tracker)
	txnIDKey := makeSharedTxnIDKey(txn.ID)
	txnDateIDKey := makeSharedTxnDateIDKey(txn)

	txnItem := SharedTxnItem{
		SK:           txnIDKey,
		GSI1SK:       txnDateIDKey,
		EntityType:   sharedTxnEntityType,
		ID:           txn.ID,
		Category:     txn.Category,
		Tracker:      txn.Tracker,
		Participants: txn.Participants,
		Date:         txn.Date,
		Amount:       txn.Amount,
		Location:     txn.Location,
		Payer:        txn.Payer,
		Details:      txn.Details,
	}

	if txn.Unsettled {
		txnItem.Unsettled = unsettledFlagTrue
	}
	// store the split map as a JSON string
	if txn.Split != nil {
		splitJSON, err := json.Marshal(txn.Split)
		if err != nil {
			return err
		}
		txnItem.SplitJSON = string(splitJSON)
	}

	for _, p := range txn.Participants {
		userIDKey := makeUserIDKey(p)
		txnItem.PK = userIDKey
		txnItem.GSI1PK = userIDKey
		err := r.table.PutItem(txnItem, option.PutCondition("attribute_exists(SK)"))
		if err != nil {
			return attrNotExistsOrErr(err)
		}
	}

	txnItem.PK = trackerIDKey
	txnItem.GSI1PK = trackerIDKey
	err := r.table.PutItem(txnItem, option.PutCondition("attribute_exists(SK)"))
	return attrNotExistsOrErr(err)
}

func (r *sharedTxnsRepo) Delete(input app.DelSharedTxnInput) error {
	trackerIDKey := makeTrackerIDKey(input.Tracker)
	txnIDKey := makeSharedTxnIDKey(input.TxnID)

	for _, p := range input.Participants {
		userIDKey := makeUserIDKey(p)
		err := r.table.DeleteItem(attributes.String(userIDKey), attributes.String(txnIDKey))
		if err != nil {
			return err
		}
	}

	return r.table.DeleteItem(attributes.String(trackerIDKey), attributes.String(txnIDKey))
}

// Settle takes a slice of SettleTxnInputs and removes the "Unsettled"
// attribute from the database items that correspond to the transactions
func (r *sharedTxnsRepo) Settle(txns []SettleTxnInput) error {
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

func makeSharedTxnDateIDKey(txn app.SharedTransaction) string {
	return fmt.Sprintf("%s#%d#%s", sharedTxnKeyPrefix, txn.Date, txn.ID)
}

func makeSharedTxnDateKey(date int64) string {
	return fmt.Sprintf("%s#%d", sharedTxnKeyPrefix, date)
}
