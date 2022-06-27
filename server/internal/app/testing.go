package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

// #region TRANSACTIONS

// NewGetTransactionRequest creates a request to be used in tests get an transaction
// by ID, with ID in the request context.
func NewGetTransactionRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/transactions/%s", id), nil)
	return req
}

// MakeTxnRequestPayload generates the payload to be given to
// NewCreateTransactionRequest
func MakeTxnRequestPayload(txn Transaction) map[string]io.Reader {
	return map[string]io.Reader{
		"location": strings.NewReader(txn.Location),
		"details":  strings.NewReader(txn.Details),
		"amount":   strings.NewReader(strconv.FormatInt(txn.Amount, 10)),
		"date":     strings.NewReader(strconv.FormatInt(txn.Date, 10)),
		"category": strings.NewReader(txn.Category),
	}
}

// NewCreateTransactionRequest creates a request to be used in tests to create a
// transaction, simulating data submitted from a form
func NewCreateTransactionRequest(values map[string]io.Reader) *http.Request {
	// TODO: refactor this and update to use similar logic
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer w.Close()
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			fw, err = w.CreateFormFile(key, x.Name())
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			// non-file values
			fw, err = w.CreateFormField(key)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			fmt.Println(err.Error())
		}
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// NewUpdateTransactionRequest creates a request to be used in tests to update a
// transaction.
func NewUpdateTransactionRequest(txn Transaction) *http.Request {
	values := MakeTxnRequestPayload(txn)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer w.Close()
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			fw, err = w.CreateFormFile(key, x.Name())
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			// non-file values
			fw, err = w.CreateFormField(key)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			fmt.Println(err.Error())
		}
	}

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/transactions/%s", txn.ID), &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// NewDeleteTransactionRequest creates a request to be used in tests to delete a
// transaction.
func NewDeleteTransactionRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/transactions/%s", id), nil)
	return req
}

// NewGetTransactionsByUserRequest creates a request to be used in tests to get all
// transactions of a user, with the user in the request context.
func NewGetTransactionsByUserRequest(userID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/transactions/user/%s", userID), nil)
	ctx := context.WithValue(req.Context(), CtxKeyUserID, userID)
	return req.WithContext(ctx)
}

// NewGetAllTransactionsRequest creates a request to be used in tests to get all
// transactions.
func NewGetAllTransactionsRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
	return req
}

// #endregion TRANSACTIONS

// #region USERS

// NewGetUserRequest creates a request to be used in tests to get a user by ID,
// with the ID in the request context.
func NewGetUserRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", id), nil)
	ctx := context.WithValue(req.Context(), CtxKeyUserID, id)
	return req.WithContext(ctx)
}

// NewCreateUserRequest creates a request to be used in tests to get create a
// new user.
func NewCreateUserRequest(userJSON []byte) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(userJSON))
	return req
}

// NewGetSelfRequest creates a request to be used in tests to get the user from
// the session.
func NewGetSelfRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/self", nil)
	return req
}

// NewGetALlUsers creates a request to be used in tests to get all users.
func NewGetAllUsersRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	return req
}

// #endregion USERS

// #region SHARED_TXNS

// NewCreateTrackerRequest creates a request to be used in tests to create a new tracker.
func NewCreateTrackerRequest(t testing.TB, trackerDetails Tracker) *http.Request {
	trackerJSON, err := json.Marshal(trackerDetails)
	if err != nil {
		t.Fatalf(err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/trackers", bytes.NewBuffer(trackerJSON))
	return req
}

// NewGetTrackerByIDRequest creates a request to be used in tests to get a
// tracker by ID, with the ID in the request context.
func NewGetTrackerByIDRequest(trackerID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/trackers/%s", trackerID), nil)
	ctx := context.WithValue(req.Context(), CtxKeyTrackerID, trackerID)
	return req.WithContext(ctx)
}

// NewGetTrackerByUserRequest creates a request to be used in tests to get a
// tracker by userID, with the userID in the request context.
func NewGetTrackerByUserRequest(userID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/trackers/user/%s", userID), nil)
	ctx := context.WithValue(req.Context(), CtxKeyUserID, userID)
	return req.WithContext(ctx)
}

// NewGetTxnsByTrackerRequest creates a request to be used in tests to get a
// a list of transactions by trackerID, with the trackerID in the request context
func NewGetTxnsByTrackerRequest(trackerID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/trackers/%s/transactions", trackerID), nil)
	ctx := context.WithValue(req.Context(), CtxKeyTrackerID, trackerID)
	return req.WithContext(ctx)
}

// MakeSharedTxnRequestPayload generates the payload to be given to NewCreateSharedTxnRequest
func MakeSharedTxnRequestPayload(txn SharedTransaction) (bytes.Buffer, string) {
	// make a comma separated list of participants
	participants := strings.Join(txn.Participants, ",")

	var unsettled string
	if txn.Unsettled {
		unsettled = "true"
	}

	values := map[string]io.Reader{
		"shop":   strings.NewReader(txn.Shop),
		"amount": strings.NewReader(strconv.FormatInt(txn.Amount, 10)),
		// NOTE: currently, date will never be empty, change this?
		"date":         strings.NewReader(strconv.FormatInt(txn.Date, 10)),
		"participants": strings.NewReader(participants),
		"unsettled":    strings.NewReader(unsettled),
		"category":     strings.NewReader(txn.Category),
		"payer":        strings.NewReader(txn.Payer),
	}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			fw, err = w.CreateFormFile(key, x.Name())
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			// non-file values
			fw, err = w.CreateFormField(key)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			fmt.Println(err.Error())
		}
	}
	w.Close()
	return b, w.FormDataContentType()
}

// NewCreateSharedTxnRequest creates a request to be used in tests to create a
// shared transaction, simulating data submitted from a form
func NewCreateSharedTxnRequest(txn SharedTransaction) *http.Request {
	b, contentType := MakeSharedTxnRequestPayload(txn)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/trackers/%s/transactions", txn.Tracker), &b)
	req.Header.Set("Content-Type", contentType)
	return req

}

// NewUpdateSharedTxnRequest creates a request to be used in tests to update a
// shared transaction
func NewUpdateSharedTxnRequest(txn SharedTransaction) *http.Request {
	b, contentType := MakeSharedTxnRequestPayload(txn)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/trackers/%s/transactions/%s", txn.Tracker, txn.ID), &b)
	req.Header.Set("Content-Type", contentType)
	return req
}

// NewGetSharedTxnByIDRequest creates a request to be used in tests to delete a
// shared transaction.
func NewDeleteSharedTxnRequest(txn SharedTransaction) *http.Request {
	input := DelSharedTxnInput{
		Tracker:      txn.Tracker,
		TxnID:        txn.ID,
		Participants: txn.Participants,
	}
	inputJSON, _ := json.Marshal(input)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/trackers/%s/transactions/%s", input.Tracker, input.TxnID), bytes.NewBuffer(inputJSON))
	return req
}

// NewGetUnsettledTxnsByTrackerRequest creates a request to be used in tests to
// get unsettled transactions
func NewGetUnsettledTxnsByTrackerRequest(trackerID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/trackers/%s/transactions/unsettled", trackerID), nil)
	ctx := context.WithValue(req.Context(), CtxKeyTrackerID, trackerID)
	return req.WithContext(ctx)
}

// NewSettleTxnsRequest creates a request to be used in tests to settle all
// transactions in a tracker
func NewSettleTxnsRequest(t testing.TB, txns []SharedTransaction) *http.Request {
	transactionsJSON, err := json.Marshal(txns)
	if err != nil {
		t.Fatalf(err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions/shared/settle", bytes.NewBuffer(transactionsJSON))
	return req
}

// #endregion SHARED_TXNS

// NewGoogleCallbackRequest creates a request to be used in tests to call the
// Google callback route.
func NewGoogleCallbackRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/callback_google", nil)
	return req
}
