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
	Date         int64    `json:"date"`
	Shop         string   `json:"shop"`
	Amount       int64    `json:"amount"`
	Category     string   `json:"category"`
	Payer        string   `json:"payer"`
	Participants []string `json:"participants"`
	Unsettled    bool     `json:"unsettled"`
	Tracker      string   `json:"tracker"`
}

// GetTxnsByTracker handles a HTTP request to get a list of transactions belonging
// to a tracker with the given ID, returning the list of transactions
func (a *App) GetTxnsByTracker(rw http.ResponseWriter, r *http.Request) {
	// TODO: should require the userID as well to check that the user is allowed to get them?
	trackerID := r.Context().Value(CtxKeyTrackerID).(string)

	transactions, err := a.store.GetTxnsByTracker(trackerID)

	if err != nil {
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
	userID := r.Context().Value(CtxKeyUserID).(string)
	// TODO: refactor, use same logic for transactions and here
	err := r.ParseMultipartForm(1024 * 1024 * 5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	participants := r.FormValue("participants")
	if participants == "" {
		http.Error(w, "participants must be provided", http.StatusBadRequest)
	}
	if !strings.Contains(participants, userID) {
		http.Error(w, "users cannot create a transaction that they are not part of", http.StatusForbidden)
		return
	}
	splitParticipants := strings.Split(participants, ",")

	shop := r.FormValue("shop")
	if shop == "" {
		http.Error(w, "transaction name must be provided", http.StatusBadRequest)
		return
	}

	amount := r.FormValue("amount")
	if amount == "0" {
		http.Error(w, "amount must be provided", http.StatusBadRequest)
		return
	}

	amountParsed, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		http.Error(w, "error parsing amount to int: "+err.Error(), http.StatusInternalServerError)
	}

	date := r.FormValue("date")
	if date == "0" {
		http.Error(w, "date must be provided", http.StatusBadRequest)
		return
	}

	dateParsed, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		http.Error(w, "error parsing date to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.store.CreateSharedTxn(SharedTransaction{
		Participants: splitParticipants,
		Shop:         shop,
		Amount:       amountParsed,
		Date:         dateParsed,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
