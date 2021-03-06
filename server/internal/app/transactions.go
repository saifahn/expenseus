package app

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
)

type Transaction struct {
	ID       string `json:"id"`
	Location string `json:"location" validate:"required"`
	UserID   string `json:"userId" validate:"required"`
	Amount   int64  `json:"amount" validate:"required"`
	Date     int64  `json:"date" validate:"required"`
	Category string `json:"category" validate:"required"`
	Details  string `json:"details"`
	ImageKey string `json:"imageKey,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
}

// GetTransaction handles a HTTP request to get an transaction by ID, returning the transaction.
func (a *App) GetTransaction(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)
	transactionID := r.Context().Value(CtxKeyTransactionID).(string)

	transaction, err := a.store.GetTransaction(userID, transactionID)

	if err != nil {
		if err == ErrDBItemNotFound {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(rw, fmt.Sprintf("something went wrong getting transaction: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	if transaction.ImageKey != "" {
		transaction, err = a.images.AddImageToTransaction(transaction)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(transaction)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTransactionsByUser handles a HTTP request to get all transactions of a user,
// returning a list of transactions.
func (a *App) GetTransactionsByUser(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)

	transactions, err := a.store.GetTransactionsByUser(userID)

	// TODO: account for different errors
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, txn := range transactions {
		if txn.ImageKey != "" {
			transactions[i], err = a.images.AddImageToTransaction(txn)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(transactions)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTxnsBetweenDates handles a HTTP request to get all transactions from a
// user between two dates given in epoch seconds, returning a list of txns.
func (a *App) GetTxnsBetweenDates(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)
	from := r.Context().Value(CtxKeyDateFrom).(int64)
	to := r.Context().Value(CtxKeyDateTo).(int64)

	transactions, err := a.store.GetTxnsBetweenDates(userID, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseTxnForm(r *http.Request, w http.ResponseWriter) *Transaction {
	err := r.ParseMultipartForm(1024 * 1024 * 5)
	if err != nil {
		if err == multipart.ErrMessageTooLarge {
			http.Error(w, "image size too large", http.StatusRequestEntityTooLarge)
			return nil
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	location := r.FormValue("location")
	details := r.FormValue("details")

	amount := r.FormValue("amount")
	amountParsed, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		http.Error(w, "error parsing amount to int: "+err.Error(), http.StatusInternalServerError)
	}

	date := r.FormValue("date")
	dateParsed, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		http.Error(w, "error parsing date to int: "+err.Error(), http.StatusInternalServerError)
		return nil
	}

	category := r.FormValue("category")

	return &Transaction{
		Location: location,
		Amount:   amountParsed,
		Date:     dateParsed,
		Category: category,
		Details:  details,
	}
}

// CreateTransaction handles a HTTP request to create a new transaction.
func (a *App) CreateTransaction(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(CtxKeyUserID).(string)
	if !ok {
		http.Error(rw, ErrUserNotInCtx.Error(), http.StatusUnauthorized)
	}

	txn := parseTxnForm(r, rw)
	txn.UserID = userID
	err := a.validate.Struct(txn)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	// don't error on missing file - it's ok not to have an image
	if err != nil && err != http.ErrMissingFile {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var imageKey string
	// upload the image only if one was supplied
	if file != nil {
		// check image is OK
		ok, err := a.images.Validate(file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(rw, "image invalid", http.StatusUnprocessableEntity)
			return
		}

		imageKey, err = a.images.Upload(file, *header)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	txn.ImageKey = imageKey

	err = a.store.CreateTransaction(*txn)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// UpdateTransaction handles a HTTP request to update a transaction.
func (a *App) UpdateTransaction(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(CtxKeyUserID).(string)
	if !ok {
		http.Error(rw, ErrUserNotInCtx.Error(), http.StatusUnauthorized)
		return
	}

	txnID, ok := r.Context().Value(CtxKeyTransactionID).(string)
	if !ok {
		http.Error(rw, ErrTxnIDNotInCtx.Error(), http.StatusBadRequest)
		return
	}

	txn := parseTxnForm(r, rw)
	txn.ID = txnID
	txn.UserID = userID
	err := a.validate.Struct(txn)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.UpdateTransaction(*txn)
	if err != nil {
		if err == ErrDBItemNotFound {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// DeleteTransaction handles a HTTP request to delete a transaction with the given ID.
func (a *App) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(CtxKeyUserID).(string)
	if !ok {
		http.Error(w, ErrUserNotInCtx.Error(), http.StatusUnauthorized)
		return
	}

	txnID, ok := r.Context().Value(CtxKeyTransactionID).(string)
	if !ok {
		http.Error(w, ErrTxnIDNotInCtx.Error(), http.StatusBadRequest)
		return
	}

	err := a.store.DeleteTransaction(txnID, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
