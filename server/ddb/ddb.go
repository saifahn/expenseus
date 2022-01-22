package ddb

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/saifahn/expenseus"
)

type dynamoDB struct {
	usersTable        UsersTable
	transactionsTable TransactionsTable
}

func New(u *UsersTable, t *TransactionsTable) *dynamoDB {
	return &dynamoDB{usersTable: *u, transactionsTable: *t}
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
	return expenseus.Expense{}, nil
}

func (d *dynamoDB) GetExpensesByUsername(id string) ([]expenseus.Expense, error) {
	return []expenseus.Expense{}, nil
}

func (d *dynamoDB) GetAllExpenses() ([]expenseus.Expense, error) {
	transactions, err := d.transactionsTable.GetAll()
	if err != nil {
		return nil, err
	}

	var expenses []expenseus.Expense
	for _, t := range transactions {
		expenses = append(expenses, expenseus.Expense{
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
