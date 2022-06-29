package ddb

import (
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
		if err == table.ErrItemNotFound {
			return app.User{}, app.ErrDBItemNotFound
		}
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

func txnToTxnItem(txn app.Transaction) TxnItem {
	userIDKey := makeUserIDKey(txn.UserID)
	transactionIDKey := makeTxnIDKey(txn.ID)
	txnDateKey := makeTxnDateIDKey(txn)

	return TxnItem{
		PK:         userIDKey,
		SK:         transactionIDKey,
		EntityType: txnEntityType,
		ID:         txn.ID,
		UserID:     txn.UserID,
		Location:   txn.Location,
		Details:    txn.Details,
		Amount:     txn.Amount,
		Date:       txn.Date,
		GSI1PK:     userIDKey,
		GSI1SK:     txnDateKey,
		Category:   txn.Category,
	}
}

func (d *dynamoDB) CreateTransaction(txn app.Transaction) error {
	transactionID := uuid.New().String()
	txn.ID = transactionID

	err := d.transactions.Create(txnToTxnItem(txn))
	if err != nil {
		return err
	}

	return nil
}

func txnItemToTxn(ti TxnItem) app.Transaction {
	return app.Transaction{
		ID:       ti.ID,
		UserID:   ti.UserID,
		Location: ti.Location,
		Details:  ti.Details,
		Amount:   ti.Amount,
		Date:     ti.Date,
		Category: ti.Category,
	}
}

func (d *dynamoDB) GetTransaction(userID, txnID string) (app.Transaction, error) {
	ti, err := d.transactions.Get(GetTxnInput{
		ID:     txnID,
		UserID: userID,
	})
	if err != nil {
		return app.Transaction{}, err
	}

	return txnItemToTxn(*ti), nil
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

func (d *dynamoDB) GetTxnsBetweenDates(userID string, from, to int64) ([]app.Transaction, error) {
	items, err := d.transactions.GetBetweenDates(userID, from, to)
	if err != nil {
		return nil, err
	}

	txns := []app.Transaction{}
	for _, ti := range items {
		txns = append(txns, txnItemToTxn(ti))
	}

	return txns, nil
}

func (d *dynamoDB) UpdateTransaction(txn app.Transaction) error {
	err := d.transactions.Update(txnToTxnItem(txn))
	if err != nil {
		if err == ErrAttrNotExists {
			return app.ErrDBItemNotFound
		}
		return err
	}

	return nil
}

func (d *dynamoDB) DeleteTransaction(txnID, user string) error {
	return d.transactions.Delete(txnID, user)
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

	trackers := []app.Tracker{}
	for _, item := range items {
		trackers = append(trackers, trackerItemToTracker(item))
	}
	return trackers, nil
}

// CreateSharedTxn calls the repository to create a new shared transaction.
func (d *dynamoDB) CreateSharedTxn(txn app.SharedTransaction) error {
	txn.ID = uuid.New().String()
	return d.sharedTxn.Create(txn)
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
		Location:     item.Location,
		Category:     item.Category,
		Payer:        item.Payer,
		Details:      item.Details,
	}
}

// GetTxnsByTracker calls the repository to get a list of transactions from a
// tracker with the given ID.
func (d *dynamoDB) GetTxnsByTracker(trackerID string) ([]app.SharedTransaction, error) {
	items, err := d.sharedTxn.GetFromTracker(trackerID)
	if err != nil {
		return nil, err
	}

	txns := []app.SharedTransaction{}
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

func (d *dynamoDB) UpdateSharedTxn(txn app.SharedTransaction) error {
	return d.sharedTxn.Update(txn)
}

// DeleteSharedTxn calls the repository to delete a shared transaction with
// the given information.
func (d *dynamoDB) DeleteSharedTxn(input app.DelSharedTxnInput) error {
	return d.sharedTxn.Delete(input)
}

// SettleTxns takes a list of transactions and calls the repository to mark
// the database items that are related as settled.
func (d *dynamoDB) SettleTxns(txns []app.SharedTransaction) error {
	var input []SettleTxnInput

	for _, t := range txns {
		input = append(input, SettleTxnInput{
			ID:           t.ID,
			TrackerID:    t.Tracker,
			Participants: t.Participants,
		})
	}

	return d.sharedTxn.Settle(input)
}
