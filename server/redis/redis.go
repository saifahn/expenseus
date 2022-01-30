package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/saifahn/expenseus/internal/app"
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

func (r *Redis) CreateExpense(ed app.ExpenseDetails) error {
	// generate id for expense
	expenseID := uuid.New().String()
	// get the time now for the score on the sets
	createdAt := time.Now().Unix()

	e := app.Expense{
		ExpenseDetails: ed,
		ID:             expenseID,
	}

	expenseJSON, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	pipe := r.db.TxPipeline()
	// record the expenseID in the expense sorted set
	pipe.ZAdd(ctx, AllExpensesKey(), &redis.Z{Score: float64(createdAt), Member: expenseID})
	// record the expenseID in the user-expense sorted set
	pipe.ZAdd(ctx, UserExpensesKey(ed.UserID), &redis.Z{Score: float64(createdAt), Member: expenseID})
	// set the expense at the expense key
	pipe.Set(ctx, ExpenseKey(expenseID), expenseJSON, 0)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetAllExpenses() ([]app.Expense, error) {
	expenseIDs := r.db.ZRange(ctx, AllExpensesKey(), 0, -1).Val()

	var expenses []app.Expense
	for _, id := range expenseIDs {
		e, err := r.GetExpense(id)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	return expenses, nil
}

func (r *Redis) GetExpensesByUsername(username string) ([]app.Expense, error) {
	// get userid from username/userid map
	userid := r.db.HGet(ctx, UsernameIDMapKey(), username).Val()
	// get expenseIDs from the user expenses set
	expenseIDs := r.db.ZRange(ctx, UserExpensesKey(userid), 0, -1).Val()

	var expenses []app.Expense
	for _, id := range expenseIDs {
		e, err := r.GetExpense(id)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	return expenses, nil
}

func (r *Redis) GetExpense(expenseID string) (app.Expense, error) {
	val, err := r.db.Get(ctx, ExpenseKey(expenseID)).Result()
	if err != nil {
		return app.Expense{}, err
	}

	var expense app.Expense
	err = json.Unmarshal([]byte(val), &expense)
	if err != nil {
		panic(err)
	}

	return expense, nil
}

func (r *Redis) CreateUser(u app.User) error {
	// get time now
	createdAt := time.Now().Unix()
	userJSON, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	pipe := r.db.TxPipeline()
	// add id to users sorted set
	pipe.ZAdd(ctx, AllUsersKey(), &redis.Z{Score: float64(createdAt), Member: u.ID})
	// add user JSON data to user:id key
	pipe.Set(ctx, UserKey(u.ID), userJSON, 0)
	// add to username:userid map
	pipe.HSet(ctx, UsernameIDMapKey(), u.Username, u.ID)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetUser(id string) (app.User, error) {
	// get the user data at the user:id key
	val, err := r.db.Get(ctx, UserKey(id)).Result()
	if err != nil {
		return app.User{}, err
	}

	// convert the user into a User struct
	var user app.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		panic(err)
	}
	return user, nil
}

func (r *Redis) GetAllUsers() ([]app.User, error) {
	// get the ids from users sorted set
	userIDs := r.db.ZRange(ctx, AllUsersKey(), 0, -1).Val()

	var users []app.User
	for _, id := range userIDs {
		user, err := r.GetUser(id)
		if err != nil {
			return []app.User{}, err
		}

		users = append(users, user)
	}
	return users, nil
}

func AllExpensesKey() string {
	return "expenses"
}

func UserExpensesKey(userid string) string {
	return fmt.Sprintf("user:%v:expenses", userid)
}

func ExpenseKey(id string) string {
	return fmt.Sprintf("expense:%v", id)
}

func AllUsersKey() string {
	return "users"
}

func UserKey(id string) string {
	return fmt.Sprintf("user:%v", id)
}

func UsernameIDMapKey() string {
	return "usernames:userids"
}
