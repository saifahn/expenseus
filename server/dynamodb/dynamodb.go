package dynamodb

import (
	"github.com/saifahn/expenseus"
)

type dynamoDB struct {
	usersTable UsersTable
}

func New(u *UsersTable) *dynamoDB {
	return &dynamoDB{usersTable: *u}
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
	return []expenseus.Expense{}, nil
}

func (d *dynamoDB) CreateExpense(ed expenseus.ExpenseDetails) error {
	return nil
}

func (d *dynamoDB) GetAllUsers() ([]expenseus.User, error) {
	return []expenseus.User{}, nil
}
