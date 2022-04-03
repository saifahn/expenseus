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

func (d *dynamoDB) GetTransaction(id string) (app.Transaction, error) {
	t, err := d.transactionsTable.Get(id)
	if err != nil {
		return app.Transaction{}, err
	}

	transaction := app.Transaction{
		ID:                 t.ID,
		TransactionDetails: t.TransactionDetails,
	}
	return transaction, nil
}

func (d *dynamoDB) GetTransactionsByUsername(username string) ([]app.Transaction, error) {
	// look up users table for user with the name
	u, err := d.usersTable.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	// then look in the transactions table for transactions with that ID
	tItems, err := d.transactionsTable.GetByUserID(u.ID)
	if err != nil {
		return nil, err
	}

	transactions := []app.Transaction{}
	for _, t := range tItems {
		transactions = append(transactions, app.Transaction{
			ID:                 t.ID,
			TransactionDetails: t.TransactionDetails,
		})
	}
	return transactions, nil
}

func (d *dynamoDB) GetAllTransactions() ([]app.Transaction, error) {
	transactionItems, err := d.transactionsTable.GetAll()
	if err != nil {
		return nil, err
	}

	var transactions []app.Transaction
	for _, t := range transactionItems {
		transactions = append(transactions, app.Transaction{
			ID:                 t.ID,
			TransactionDetails: t.TransactionDetails,
		})
	}

	return transactions, nil
}
