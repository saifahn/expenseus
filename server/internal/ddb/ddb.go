package ddb

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/saifahn/expenseus/internal/app"
)

type dynamoDB struct {
	users        UserRepository
	transactions TxnRepository
	trackers     TrackersRepository
	sharedTxn    SharedTxnsRepository
}

func New(d dynamodbiface.DynamoDBAPI, tableName string) *dynamoDB {
	tbl := table.New(d, tableName)
	users := NewUserRepository(tbl)
	transactions := NewTxnRepository(tbl)
	trackers := NewTrackersRepository(tbl)
	sharedTxn := NewSharedTxnsRepository(tbl)

	return &dynamoDB{users, transactions, trackers, sharedTxn}
}

func userToUserItem(u app.User) UserItem {
	userIDKey := makeUserIDKey(u.ID)
	return UserItem{
		PK:         userIDKey,
		SK:         userIDKey,
		EntityType: userEntityType,
		ID:         u.ID,
		Username:   u.Username,
		Name:       u.Name,
		GSI1PK:     allUsersKey,
		GSI1SK:     userIDKey,
	}
}

func (d *dynamoDB) CreateUser(u app.User) error {
	err := d.users.PutIfNotExists(userToUserItem(u))
	if err != nil {
		return err
	}

	return nil
}

func userItemToUser(ui UserItem) app.User {
	return app.User{
		ID:       ui.ID,
		Username: ui.Username,
		Name:     ui.Name,
	}
}

func (d *dynamoDB) GetUser(id string) (app.User, error) {
	ui, err := d.users.Get(id)
	if err != nil {
		return app.User{}, err
	}

	return userItemToUser(ui), nil
}

func (d *dynamoDB) GetAllUsers() ([]app.User, error) {
	userItems, err := d.users.GetAll()
	if err != nil {
		return nil, err
	}

	var users []app.User
	for _, ui := range userItems {
		users = append(users, userItemToUser(ui))
	}
	return users, nil
}

func (d *dynamoDB) CreateTransaction(td app.TransactionDetails) error {
	userIDKey := makeUserIDKey(td.UserID)
	// generate an ID for the transaction
	transactionID := uuid.New().String()
	transactionIDKey := makeTxnIDKey(transactionID)

	item := &TransactionItem{
		PK:         userIDKey,
		SK:         transactionIDKey,
		EntityType: txnEntityType,
		ID:         transactionID,
		UserID:     td.UserID,
		Name:       td.Name,
		Amount:     td.Amount,
		Date:       td.Date,
		GSI1PK:     allTxnKey,
		GSI1SK:     transactionIDKey,
	}
	err := d.transactions.PutIfNotExists(*item)
	if err != nil {
		return err
	}

	log.Println("transaction successfully created")
	return nil
}

func txnItemToTxn(ti TransactionItem) app.Transaction {
	return app.Transaction{
		ID: ti.ID,
		TransactionDetails: app.TransactionDetails{
			UserID: ti.UserID,
			Name:   ti.Name,
			Amount: ti.Amount,
			Date:   ti.Date,
		},
	}
}

func (d *dynamoDB) GetTransaction(transactionID string) (app.Transaction, error) {
	ti, err := d.transactions.Get(transactionID)
	if err != nil {
		return app.Transaction{}, err
	}

	return txnItemToTxn(*ti), nil
}

func (d *dynamoDB) GetAllTransactions() ([]app.Transaction, error) {
	transactionItems, err := d.transactions.GetAll()
	if err != nil {
		return nil, err
	}

	var transactions []app.Transaction
	for _, ti := range transactionItems {
		transactions = append(transactions, txnItemToTxn(ti))
	}

	return transactions, nil
}

func (d *dynamoDB) GetTransactionsByUser(userID string) ([]app.Transaction, error) {
	items, err := d.transactions.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	transactions := []app.Transaction{}
	for _, ti := range items {
		transactions = append(transactions, txnItemToTxn(ti))
	}
	return transactions, nil
}

// CreateTracker calls the repository to create a new tracker item.
func (d *dynamoDB) CreateTracker(tracker app.Tracker) error {
	id := uuid.New().String()

	err := d.trackers.Create(CreateTrackerInput{
		ID:    id,
		Name:  tracker.Name,
		Users: tracker.Users,
	})
	if err != nil {
		return err
	}

	return nil
}

func trackerItemToTracker(ti TrackerItem) app.Tracker {
	return app.Tracker{
		ID:    ti.ID,
		Name:  ti.Name,
		Users: ti.Users,
	}
}

// GetTracker calls the repository to get a tracker by its ID, returning the
// tracker if found and an error if not.
func (d *dynamoDB) GetTracker(trackerID string) (app.Tracker, error) {
	item, err := d.trackers.Get(trackerID)
	if err != nil {
		return app.Tracker{}, err
	}

	return trackerItemToTracker(*item), nil
}

// GetTrackersByUser calls the repository to return a list of trackers that a
// user belongs to.
func (d *dynamoDB) GetTrackersByUser(userID string) ([]app.Tracker, error) {
	items, err := d.trackers.GetByUser(userID)
	if err != nil {
		return nil, err
	}

	var trackers []app.Tracker
	for _, item := range items {
		trackers = append(trackers, trackerItemToTracker(item))
	}
	return trackers, nil
}

// CreateSharedTxn calls the repository to create a new shared transaction.
func (d *dynamoDB) CreateSharedTxn(txn app.SharedTransaction) error {
	// TODO: create the ID
	// TODO: also just pass in the transaction, don't need this input thing any more
	return d.sharedTxn.Create(CreateSharedTxnInput{
		ID:           txn.ID,
		TrackerID:    txn.Tracker,
		Participants: txn.Participants,
		Unsettled:    txn.Unsettled,
		Date:         txn.Date,
		Amount:       txn.Amount,
		Shop:         txn.Shop,
	})
}

// A helper function for converting an item in the database representing a
// shared transaction to a shared transaction structure for the application.
func sharedTxnItemToSharedTxn(item SharedTxnItem) app.SharedTransaction {
	return app.SharedTransaction{
		ID:           item.ID,
		Participants: item.Participants,
		Date:         item.Date,
		Amount:       item.Amount,
		Tracker:      item.Tracker,
		Unsettled:    item.Unsettled == unsettledFlagTrue,
		Shop:         item.Shop,
	}
}

// GetTxnsByTracker calls the repository to get a list of transactions from a
// tracker with the given ID.
func (d *dynamoDB) GetTxnsByTracker(trackerID string) ([]app.SharedTransaction, error) {
	items, err := d.sharedTxn.GetFromTracker(trackerID)
	if err != nil {
		return nil, err
	}

	var txns []app.SharedTransaction
	for _, i := range items {
		txns = append(txns, sharedTxnItemToSharedTxn(i))
	}
	return txns, nil
}

// GetUnsettledTxnsByTracker calls the repository to get a list of unsettled
// transactions from a tracker with the given ID.
func (d *dynamoDB) GetUnsettledTxnsByTracker(trackerID string) ([]app.SharedTransaction, error) {
	items, err := d.sharedTxn.GetUnsettledFromTracker(trackerID)
	if err != nil {
		return nil, err
	}

	var txns []app.SharedTransaction
	for _, i := range items {
		txns = append(txns, sharedTxnItemToSharedTxn(i))
	}
	return txns, nil
}

// SettleTxns takes a list of transactions and calls the repository to mark
// the database items that are related as settled.
func (d *dynamoDB) SettleTxns(txns []app.SharedTransaction) error {
	var input []SettleTxnInputItem

	for _, t := range txns {
		input = append(input, SettleTxnInputItem{
			ID:           t.ID,
			TrackerID:    t.Tracker,
			Participants: t.Participants,
		})
	}

	return d.sharedTxn.Settle(input)
}
