package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SharedTransaction struct {
	ID           string   `json:"id"`
	Date         int64    `json:"date" validate:"required"`
	Shop         string   `json:"shop" validate:"required"`
	Amount       int64    `json:"amount" validate:"required"`
	Category     string   `json:"category" validate:"required"`
	Payer        string   `json:"payer"`
	Participants []string `json:"participants" validate:"required,min=1"`
	Unsettled    bool     `json:"unsettled"`
	Tracker      string   `json:"tracker" validate:"required"`
}

// GetTxnsByTracker handles a HTTP request to get a list of transactions belonging
// to a tracker with the given ID, returning the list of transactions
func (a *App) GetTxnsByTracker(rw http.ResponseWriter, r *http.Request) {
	// TODO: should require the userID as well to check that the user is allowed to get them?
	trackerID := r.Context().Value(CtxKeyTrackerID).(string)

	transactions, err := a.store.GetTxnsByTracker(trackerID)

	if err != nil {
		if err == ErrDBItemNotFound {
			http.Error(rw, "a tracker with the given trackerID does not exist", http.StatusNotFound)
			return
		}
		http.Error(rw, fmt.Sprintf("something went wrong getting shared transactions from tracker: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(transactions)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateSharedTxn handles a HTTP request to create a shared transaction
func (a *App) CreateSharedTxn(w http.ResponseWriter, r *http.Request) {
	tracker := r.Context().Value(CtxKeyTrackerID).(string)
	userID := r.Context().Value(CtxKeyUserID).(string)
	// TODO: refactor, use same logic for transactions and here
	err := r.ParseMultipartForm(1024 * 1024 * 5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	participants := r.FormValue("participants")
	if !strings.Contains(participants, userID) {
		http.Error(w, "users cannot create a transaction that they are not part of", http.StatusForbidden)
		return
	}
	splitParticipants := strings.Split(participants, ",")

	amount := r.FormValue("amount")
	amountParsed, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		http.Error(w, "error parsing amount to int: "+err.Error(), http.StatusInternalServerError)
	}

	date := r.FormValue("date")
	dateParsed, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		http.Error(w, "error parsing date to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	unsettled := r.FormValue("unsettled") == "true"
	shop := r.FormValue("shop")
	category := r.FormValue("category")
	payer := r.FormValue("payer")

	txn := SharedTransaction{
		Participants: splitParticipants,
		Shop:         shop,
		Amount:       amountParsed,
		Date:         dateParsed,
		Tracker:      tracker,
		Unsettled:    unsettled,
		Category:     category,
		Payer:        payer,
	}
	err = a.validate.Struct(txn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.CreateSharedTxn(txn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// GetUnsettledTxnsByTracker handles a HTTP request to get transactions that
// are unsettled and returns the list of transactions
func (a *App) GetUnsettledTxnsByTracker(w http.ResponseWriter, r *http.Request) {
	trackerID := r.Context().Value(CtxKeyTrackerID).(string)

	transactions, err := a.store.GetUnsettledTxnsByTracker(trackerID)
	if err != nil {
		if err == ErrDBItemNotFound {
			http.Error(w, "a tracker with the given trackerID does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("something went wrong getting unsettled shared transactions from tracker: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SettleTxns handles a HTTP request to mark all transactions in a tracker as
// settled
func (a *App) SettleTxns(w http.ResponseWriter, r *http.Request) {
	// TODO: make sure the user is part of the tracker
	var txns []SharedTransaction
	err := json.NewDecoder(r.Body).Decode(&txns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.SettleTxns(txns)
	if err != nil {
		http.Error(w, fmt.Sprintf("something went wrong settling transactions for tracker: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
