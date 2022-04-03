package app

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
)

const CtxKeyTransactionID contextKey = iota

type TransactionDetails struct {
	Name     string `json:"name"`
	UserID   string `json:"userId"`
	Amount   int64  `json:"amount"`
	Date     int64  `json:"date"`
	ImageKey string `json:"imageKey,omitempty"`
}

type Transaction struct {
	TransactionDetails
	ID       string `json:"id"`
	ImageURL string `json:"imageUrl,omitempty"`
}

// GetTransaction handles a HTTP request to get an transaction by ID, returning the transaction.
func (a *App) GetTransaction(rw http.ResponseWriter, r *http.Request) {
	transactionID := r.Context().Value(CtxKeyTransactionID).(string)

	transaction, err := a.store.GetTransaction(transactionID, transactionID)

	// TODO: should account for different kinds of errors
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
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
	username := r.Context().Value(CtxKeyUsername).(string)

	transactions, err := a.store.GetTransactionsByUser(username)

	// TODO: account for different errors
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(transactions)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAllTransactions handles a HTTP request to get all transactions, returning a list
// of transactions.
func (a *App) GetAllTransactions(rw http.ResponseWriter, r *http.Request) {
	transactions, err := a.store.GetAllTransactions()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, e := range transactions {
		if e.ImageKey != "" {
			transactions[i], err = a.images.AddImageToTransaction(e)
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

// CreateTransaction handles a HTTP request to create a new transaction.
func (a *App) CreateTransaction(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(CtxKeyUserID).(string)
	if !ok {
		http.Error(rw, "user id not found in context", http.StatusUnauthorized)
	}

	err := r.ParseMultipartForm(1024 * 1024 * 5)
	if err != nil {
		if err == multipart.ErrMessageTooLarge {
			http.Error(rw, "image size too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	transactionName := r.FormValue("transactionName")
	if transactionName == "" {
		http.Error(rw, "transaction name not found", http.StatusBadRequest)
		return
	}

	amount := r.FormValue("amount")
	if amount == "" {
		http.Error(rw, "amount not present", http.StatusBadRequest)
		return
	}

	amountParsed, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		http.Error(rw, "error parsing amount to int: "+err.Error(), http.StatusInternalServerError)
	}

	date := r.FormValue("date")
	if date == "" {
		http.Error(rw, "date not present", http.StatusBadRequest)
		return
	}

	dateParsed, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		http.Error(rw, "error parsing date to int: "+err.Error(), http.StatusInternalServerError)
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

	err = a.store.CreateTransaction(TransactionDetails{
		Name:     transactionName,
		UserID:   userID,
		ImageKey: imageKey,
		Amount:   amountParsed,
		Date:     dateParsed,
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}
