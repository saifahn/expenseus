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
}

func New(d dynamodbiface.DynamoDBAPI, tableName string) *dynamoDB {
	users := NewUserRepository(table.New(d, tableName))
	transactions := NewTxnRepository(table.New(d, tableName))

	return &dynamoDB{users: users, transactions: transactions}
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

func (d *dynamoDB) GetTransaction(userID, transactionID string) (app.Transaction, error) {
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
