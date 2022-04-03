package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/saifahn/expenseus/internal/app"
)

type dynamoDB struct {
	usersTable        UsersTable
	transactionsTable TransactionsTable
}

func New(d dynamodbiface.DynamoDBAPI, usersTableName, transactionsTableName string) *dynamoDB {
	uTbl := table.New(d, usersTableName)
	usersTable := NewUsersTable(uTbl)
	tTbl := table.New(d, transactionsTableName)
	transactionsTable := NewTransactionsTable(tTbl)

	return &dynamoDB{usersTable: usersTable, transactionsTable: transactionsTable}
}

func userToUserItem(u app.User) UserItem {
	userIDKey := fmt.Sprintf("%s#%s", UserKeyPrefix, u.ID)
	return UserItem{
		PK:         userIDKey,
		SK:         userIDKey,
		EntityType: UserEntityType,
		ID:         u.ID,
		GSI1PK:     allUsersKey,
		GSI1SK:     userIDKey,
	}
}

func (d *dynamoDB) CreateUser(u app.User) error {
	err := d.usersTable.PutIfNotExists(userToUserItem(u))
	if err != nil {
		return err
	}

	return nil
}

func userItemToUser(ui UserItem) app.User {
	return app.User{
		ID: ui.ID,
	}
}

func (d *dynamoDB) GetUser(id string) (app.User, error) {
	ui, err := d.usersTable.Get(id)
	if err != nil {
		return app.User{}, err
	}

	return userItemToUser(ui), nil
}

func (d *dynamoDB) GetAllUsers() ([]app.User, error) {
	userItems, err := d.usersTable.GetAll()
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
	// generate an ID
	transactionID := uuid.New().String()
	userIDKey := fmt.Sprintf("%s#%s", UserKeyPrefix, td.UserID)
	transactionIDKey := fmt.Sprintf("%s#%s", TransactionKeyPrefix, transactionID)

	item := &TransactionItem{
		PK:         userIDKey,
		SK:         transactionIDKey,
		EntityType: transactionEntityType,
		ID:         transactionID,
		UserID:     td.UserID,
		GSI1PK:     allTxnKey,
		GSI1SK:     transactionIDKey,
	}
	err := d.transactionsTable.PutIfNotExists(*item)
	if err != nil {
		return err
	}

	fmt.Println("transaction successfully created")
	return nil
}

func transactionItemToTransaction(ti TransactionItem) app.Transaction {
	return app.Transaction{
		ID: ti.ID,
		TransactionDetails: app.TransactionDetails{
			UserID: ti.UserID,
		},
	}
}

func (d *dynamoDB) GetTransaction(userID, transactionID string) (app.Transaction, error) {
	ti, err := d.transactionsTable.Get(userID, transactionID)
	if err != nil {
		return app.Transaction{}, err
	}

	return transactionItemToTransaction(*ti), nil
}

func (d *dynamoDB) GetAllTransactions() ([]app.Transaction, error) {
	transactionItems, err := d.transactionsTable.GetAll()
	if err != nil {
		return nil, err
	}

	var transactions []app.Transaction
	for _, ti := range transactionItems {
		transactions = append(transactions, transactionItemToTransaction(ti))
	}

	return transactions, nil
}

func (d *dynamoDB) GetTransactionsByUser(userID string) ([]app.Transaction, error) {
	items, err := d.transactionsTable.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	transactions := []app.Transaction{}
	for _, ti := range items {
		transactions = append(transactions, transactionItemToTransaction(ti))
	}
	return transactions, nil
}
