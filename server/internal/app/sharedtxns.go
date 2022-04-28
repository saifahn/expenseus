package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SharedTransaction struct {
	ID           string   `json:"id"`
	Date         int64    `json:"date"`
	Shop         string   `json:"shop"`
	Amount       string   `json:"amount"`
	Category     string   `json:"category"`
	Payer        string   `json:"payer"`
	Participants []string `json:"participants"`
	Unsettled    bool     `json:"unsettled"`
	Tracker      string   `json:"tracker"`
}

// GetTxnsByTracker handles a HTTP request to get a list of transactions belonging
// to a tracker with the given ID, returning the list of transactions
func (a *App) GetTxnsByTracker(rw http.ResponseWriter, r *http.Request) {
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
