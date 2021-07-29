package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/saifahn/expenseus"
)

var ctx = context.Background()

func InitClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return client
}

func New(address string) *Redis {
	client := InitClient(address)
	return &Redis{db: *client}
}

type Redis struct {
	db redis.Client
}

func (r *Redis) RecordExpense(e expenseus.Expense) error {
	// generate id for expense
	expenseID := uuid.New().String()
	// get the time now for the score on the sets
	createdAt := time.Now().Unix()
	// get the user for the user-expense set
	user := e.User

	expense, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}

	pipe := r.db.TxPipeline()
	// record the expenseID in the expense sorted set
	pipe.ZAdd(ctx, AllExpensesKey(), &redis.Z{Score: float64(createdAt), Member: expenseID})
	// record the expenseID in the user-expense sorted set
	pipe.ZAdd(ctx, UserExpensesKey(user), &redis.Z{Score: float64(createdAt), Member: expenseID})
	// set the expense at the expense key
	pipe.Set(ctx, ExpenseKey(expenseID), expense, 0)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetAllExpenses() ([]expenseus.Expense, error) {
	expenseIDs := r.db.ZRange(ctx, AllExpensesKey(), 0, -1).Val()

	var expenses []expenseus.Expense
	for _, id := range expenseIDs {
		// TODO: handle this error
		val, err := r.db.Get(ctx, ExpenseKey(id)).Result()

		var expense expenseus.Expense
		err = json.Unmarshal([]byte(val), &expense)
		if err != nil {
			panic(err)
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (r *Redis) GetExpensesByUser(username string) ([]expenseus.Expense, error) {
	return []expenseus.Expense{}, nil
}

func (r *Redis) GetExpense(expenseID string) (expenseus.Expense, error) {
	return expenseus.Expense{}, nil
}

func AllExpensesKey() string {
	return "expenses"
}

func UserExpensesKey(user string) string {
	return fmt.Sprintf("user:%v:expenses", user)
}

func ExpenseKey(id string) string {
	return fmt.Sprintf("expense:%v", id)
}
