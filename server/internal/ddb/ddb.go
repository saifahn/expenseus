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

func (d *dynamoDB) CreateUser(u app.User) error {
	item := &UserItem{
		User: u,
	}
	err := d.usersTable.PutIfNotExists(*item)
	if err != nil {
		return err
	}

	return nil
}

func (d *dynamoDB) GetUser(id string) (app.User, error) {
	u, err := d.usersTable.Get(id)
	if err != nil {
		return app.User{}, err
	}
	user := u.User

	return user, nil
}

func (d *dynamoDB) GetExpense(id string) (app.Transaction, error) {
	t, err := d.transactionsTable.Get(id)
	if err != nil {
		return app.Transaction{}, err
	}

	expense := app.Transaction{
		ID:                 t.ID,
		TransactionDetails: t.TransactionDetails,
	}
	return expense, nil
}

func (d *dynamoDB) GetExpensesByUsername(username string) ([]app.Transaction, error) {
	// look up users table for user with the name
	u, err := d.usersTable.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	// then look in the transactions table for expenses with that ID
	tItems, err := d.transactionsTable.GetByUserID(u.ID)
	if err != nil {
		return nil, err
	}

	expenses := []app.Transaction{}
	for _, t := range tItems {
		expenses = append(expenses, app.Transaction{
			ID:                 t.ID,
			TransactionDetails: t.TransactionDetails,
		})
	}
	return expenses, nil
}

func (d *dynamoDB) GetAllExpenses() ([]app.Transaction, error) {
	transactions, err := d.transactionsTable.GetAll()
	if err != nil {
		return nil, err
	}

	var expenses []app.Transaction
	for _, t := range transactions {
		expenses = append(expenses, app.Transaction{
			ID:                 t.ID,
			TransactionDetails: t.TransactionDetails,
		})
	}

	return expenses, nil
}

func (d *dynamoDB) CreateExpense(ed app.TransactionDetails) error {
	// generate an ID
	expenseID := uuid.New().String()
	item := &TransactionItem{
		ID:                 expenseID,
		TransactionDetails: ed,
	}
	err := d.transactionsTable.PutIfNotExists(*item)
	if err != nil {
		return err
	}

	fmt.Println("expense successfully created")
	return nil
}

func (d *dynamoDB) GetAllUsers() ([]app.User, error) {
	userItems, err := d.usersTable.GetAll()
	if err != nil {
		return nil, err
	}

	var users []app.User
	for _, u := range userItems {
		users = append(users, u.User)
	}
	return users, nil
}
