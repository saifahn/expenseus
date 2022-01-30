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

func (r *Redis) CreateTransaction(ed app.TransactionDetails) error {
	// generate id for transaction
	transactionID := uuid.New().String()
	// get the time now for the score on the sets
	createdAt := time.Now().Unix()

	e := app.Transaction{
		TransactionDetails: ed,
		ID:                 transactionID,
	}

	transactionJSON, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	pipe := r.db.TxPipeline()
	// record the transactionID in the transaction sorted set
	pipe.ZAdd(ctx, AllTransactionsKey(), &redis.Z{Score: float64(createdAt), Member: transactionID})
	// record the transactionID in the user-transaction sorted set
	pipe.ZAdd(ctx, UserTransactionsKey(ed.UserID), &redis.Z{Score: float64(createdAt), Member: transactionID})
	// set the transaction at the transaction key
	pipe.Set(ctx, TransactionKey(transactionID), transactionJSON, 0)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetAllTransactions() ([]app.Transaction, error) {
	transactionIDs := r.db.ZRange(ctx, AllTransactionsKey(), 0, -1).Val()

	var transactions []app.Transaction
	for _, id := range transactionIDs {
		e, err := r.GetTransaction(id)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, e)
	}

	return transactions, nil
}

func (r *Redis) GetTransactionsByUsername(username string) ([]app.Transaction, error) {
	// get userid from username/userid map
	userid := r.db.HGet(ctx, UsernameIDMapKey(), username).Val()
	// get transactionIDs from the user transactions set
	transactionIDs := r.db.ZRange(ctx, UserTransactionsKey(userid), 0, -1).Val()

	var transactions []app.Transaction
	for _, id := range transactionIDs {
		e, err := r.GetTransaction(id)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, e)
	}

	return transactions, nil
}

func (r *Redis) GetTransaction(transactionID string) (app.Transaction, error) {
	val, err := r.db.Get(ctx, TransactionKey(transactionID)).Result()
	if err != nil {
		return app.Transaction{}, err
	}

	var transaction app.Transaction
	err = json.Unmarshal([]byte(val), &transaction)
	if err != nil {
		panic(err)
	}

	return transaction, nil
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

func AllTransactionsKey() string {
	return "transactions"
}

func UserTransactionsKey(userid string) string {
	return fmt.Sprintf("user:%v:transactions", userid)
}

func TransactionKey(id string) string {
	return fmt.Sprintf("transaction:%v", id)
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
