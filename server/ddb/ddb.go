package ddb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/nabeken/aws-go-dynamodb/table"
	"github.com/saifahn/expenseus"
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

func (d *dynamoDB) CreateUser(u expenseus.User) error {
	item := &UserItem{
		User: u,
	}
	err := d.usersTable.PutIfNotExists(*item)
	if err != nil {
		return err
	}

	return nil
}

func (d *dynamoDB) GetUser(id string) (expenseus.User, error) {
	u, err := d.usersTable.Get(id)
	if err != nil {
		return expenseus.User{}, err
	}
	user := u.User

	return user, nil
}

func (d *dynamoDB) GetExpense(id string) (expenseus.Expense, error) {
	t, err := d.transactionsTable.Get(id)
	if err != nil {
		return expenseus.Expense{}, err
	}

	expense := expenseus.Expense{
		ID:             t.ID,
		ExpenseDetails: t.ExpenseDetails,
	}
	return expense, nil
}

func (d *dynamoDB) GetExpensesByUsername(username string) ([]expenseus.Expense, error) {
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

	expenses := []expenseus.Expense{}
	for _, t := range tItems {
		expenses = append(expenses, expenseus.Expense{
			ID:             t.ID,
			ExpenseDetails: t.ExpenseDetails,
		})
	}
	return expenses, nil
}

func (d *dynamoDB) GetAllExpenses() ([]expenseus.Expense, error) {
	transactions, err := d.transactionsTable.GetAll()
	if err != nil {
		return nil, err
	}

	var expenses []expenseus.Expense
	for _, t := range transactions {
		expenses = append(expenses, expenseus.Expense{
			ID:             t.ID,
			ExpenseDetails: t.ExpenseDetails,
		})
	}

	return expenses, nil
}

func (d *dynamoDB) CreateExpense(ed expenseus.ExpenseDetails) error {
	// generate an ID
	expenseID := uuid.New().String()
	item := &TransactionItem{
		ID:             expenseID,
		ExpenseDetails: ed,
	}
	err := d.transactionsTable.PutIfNotExists(*item)
	if err != nil {
		return err
	}

	fmt.Println("expense successfully created")
	return nil
}

func (d *dynamoDB) GetAllUsers() ([]expenseus.User, error) {
	userItems, err := d.usersTable.GetAll()
	if err != nil {
		return nil, err
	}

	var users []expenseus.User
	for _, u := range userItems {
		users = append(users, u.User)
	}
	return users, nil
}
