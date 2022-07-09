package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SharedTransaction struct {
	ID           string             `json:"id"`
	Date         int64              `json:"date" validate:"required"`
	Location     string             `json:"location" validate:"required"`
	Amount       int64              `json:"amount" validate:"required"`
	Category     string             `json:"category" validate:"required"`
	Payer        string             `json:"payer" validate:"required"`
	Participants []string           `json:"participants" validate:"required,min=1"`
	Unsettled    bool               `json:"unsettled"`
	Tracker      string             `json:"tracker" validate:"required"`
	Details      string             `json:"details"`
	Split        map[string]float64 `json:"split"`
}

// GetTxnsByTracker handles a HTTP request to get a list of transactions belonging
// to a tracker with the given ID, returning the list of transactions
func (a *App) GetTxnsByTracker(rw http.ResponseWriter, r *http.Request) {
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

// GetTxnsByTrackerBetweenDates handles a HTTP request to get a list of txns
// from a tracker with the given ID and between two given dates, returning the
// list of txns
func (a *App) GetTxnsByTrackerBetweenDates(w http.ResponseWriter, r *http.Request) {
	trackerID := r.Context().Value(CtxKeyTrackerID).(string)
	from := r.Context().Value(CtxKeyDateFrom).(int64)
	to := r.Context().Value(CtxKeyDateTo).(int64)

	txns, err := a.store.GetTxnsByTrackerBetweenDates(trackerID, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(w).Encode(txns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseSharedTxnForm(r *http.Request, w http.ResponseWriter) *SharedTransaction {
	err := r.ParseMultipartForm(1000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	amount := r.FormValue("amount")
	amountParsed, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		http.Error(w, "error parsing amount to int: "+err.Error(), http.StatusInternalServerError)
		return nil
	}

	date := r.FormValue("date")
	dateParsed, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		http.Error(w, "error parsing date to int: "+err.Error(), http.StatusInternalServerError)
		return nil
	}

	unsettled := r.FormValue("unsettled") == "true"
	location := r.FormValue("location")
	category := r.FormValue("category")
	payer := r.FormValue("payer")
	details := r.FormValue("details")

	txn := &SharedTransaction{
		Location:  location,
		Amount:    amountParsed,
		Date:      dateParsed,
		Unsettled: unsettled,
		Category:  category,
		Payer:     payer,
		Details:   details,
	}

	split := r.FormValue("split")
	if split != "" {
		splitMap := map[string]float64{}
		// the format will be "userid:split,userid:split", so split them into individual parts
		userSplits := strings.Split(split, ",")
		for _, u := range userSplits {
			// split into [userid, split], then assign to the map
			userSplitSeparated := strings.Split(u, ":")
			splitFloat, err := strconv.ParseFloat(userSplitSeparated[1], 64)
			if err != nil {
				http.Error(w, "error parsing split into float: "+err.Error(), http.StatusInternalServerError)
			}
			splitMap[userSplitSeparated[0]] = splitFloat
		}
		txn.Split = splitMap
	}

	return txn
}

// CreateSharedTxn handles a HTTP request to create a shared transaction
func (a *App) CreateSharedTxn(w http.ResponseWriter, r *http.Request) {
	tracker := r.Context().Value(CtxKeyTrackerID).(string)
	userID := r.Context().Value(CtxKeyUserID).(string)

	participants := r.FormValue("participants")
	if !strings.Contains(participants, userID) {
		http.Error(w, "users cannot create a transaction that they are not part of", http.StatusForbidden)
		return
	}
	splitParticipants := strings.Split(participants, ",")

	txn := parseSharedTxnForm(r, w)
	txn.Participants = splitParticipants
	txn.Tracker = tracker

	err := a.validate.Struct(txn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.CreateSharedTxn(*txn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// UpdateSharedTxn handles a HTTP request to update a shared transaction
func (a *App) UpdateSharedTxn(w http.ResponseWriter, r *http.Request) {
	tracker := r.Context().Value(CtxKeyTrackerID).(string)
	userID := r.Context().Value(CtxKeyUserID).(string)
	txnID := r.Context().Value(CtxKeyTransactionID).(string)

	participants := r.FormValue("participants")
	if !strings.Contains(participants, userID) {
		http.Error(w, "users cannot create a transaction that they are not part of", http.StatusForbidden)
		return
	}
	splitParticipants := strings.Split(participants, ",")

	txn := parseSharedTxnForm(r, w)
	txn.ID = txnID
	txn.Participants = splitParticipants
	txn.Tracker = tracker

	err := a.validate.Struct(txn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.UpdateSharedTxn(*txn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type DelSharedTxnInput struct {
	Tracker      string   `json:"tracker" validate:"required"`
	TxnID        string   `json:"txnID" validate:"required"`
	Participants []string `json:"participants" validate:"required,min=1"`
}

// DeleteSharedTxn handles a HTTP request to delete a shared transaction
func (a *App) DeleteSharedTxn(w http.ResponseWriter, r *http.Request) {
	// TODO: check that the context matches the input
	// tracker := r.Context().Value(CtxKeyTrackerID).(string)
	// txnID := r.Context().Value(CtxKeyTransactionID).(string)
	var input DelSharedTxnInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.validate.Struct(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.DeleteSharedTxn(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type UnsettledResponse struct {
	Txns       []SharedTransaction `json:"transactions"`
	Debtor     string              `json:"debtor"`
	Debtee     string              `json:"debtee"`
	AmountOwed float64             `json:"amountOwed"`
}

func CalculateDebts(currentUser string, txns []SharedTransaction) UnsettledResponse {
	defaultSplit := 0.5
	var otherUser string
	var total float64

	for _, t := range txns {
		// set the other user only once as all transactions should come from the
		// same tracker and have the same participants
		if otherUser == "" {
			for _, u := range t.Participants {
				if u != currentUser {
					otherUser = u
					continue
				}
			}
		}

		var split float64
		if t.Split == nil {
			split = defaultSplit
		}

		// calculate from the perspective that the logged in user is the one who has paid
		currentUserIsPayer := t.Payer == currentUser
		if currentUserIsPayer {
			// Split represents the proportion each participant will pay for a purchase
			// so when calculating the debt, it is the inverse proportion (and is equal to
			// the other person's proportion) that is used
			if split != defaultSplit {
				split = t.Split[otherUser]
			}
			total += float64(t.Amount) * split
		} else {
			if split != defaultSplit {
				split = t.Split[currentUser]
			}
			total -= float64(t.Amount) * split
		}
	}
	return UnsettledResponse{
		Txns:       txns,
		Debtor:     otherUser,
		Debtee:     currentUser,
		AmountOwed: total,
	}
}

// GetUnsettledTxnsByTracker handles a HTTP request to get transactions that
// are unsettled and returns the list of transactions
func (a *App) GetUnsettledTxnsByTracker(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)
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
	totals := CalculateDebts(userID, transactions)

	w.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(w).Encode(totals)
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
