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

var (
	TestSeanUser = User{
		Username: "saifahn",
		Name:     "Sean Li",
		ID:       "sean_id",
	}
	TestTomomiUser = User{
		Username: "tomochi",
		Name:     "Tomomi Kinoshita",
		ID:       "tomomi_id",
	}

	TestSeanTransactionDetails = TransactionDetails{
		Name:   "Transaction 1",
		UserID: TestSeanUser.ID,
		Amount: 123,
		Date:   1644085875,
	}
	TestSeanTransaction = Transaction{
		ID:                 "1",
		TransactionDetails: TestSeanTransactionDetails,
	}

	TestTomomiTransactionDetails = TransactionDetails{
		Name:   "Transaction 2",
		UserID: TestTomomiUser.ID,
		Amount: 456,
		Date:   1644085876,
	}
	TestTomomiTransaction = Transaction{
		ID:                 "2",
		TransactionDetails: TestTomomiTransactionDetails,
	}

	TestTomomiTransaction2Details = TransactionDetails{
		Name:   "Transaction 3",
		UserID: TestTomomiUser.ID,
		Amount: 789,
		Date:   1644085877,
	}
	TestTomomiTransaction2 = Transaction{
		ID:                 "3",
		TransactionDetails: TestTomomiTransaction2Details,
	}

	TestTransactionWithImage = Transaction{
		ID: "123",
		TransactionDetails: TransactionDetails{
			Name:     "TransactionWithImage",
			UserID:   "an_ID",
			ImageKey: "test-image-key",
		},
	}

	TestTracker = Tracker{
		Name:  "Test Tracker",
		Users: []string{TestSeanUser.ID},
		ID:    "test-id",
	}
)

// NewGetTransactionRequest creates a request to be used in tests get an transaction
// by ID, with ID in the request context.
func NewGetTransactionRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/transactions/%s", id), nil)
	ctx := context.WithValue(req.Context(), CtxKeyTransactionID, id)
	return req.WithContext(ctx)
}

// MakeCreateTransactionRequestPayload generates the payload to be given to
// NewCreateTransactionRequest
func MakeCreateTransactionRequestPayload(td TransactionDetails) map[string]io.Reader {
	return map[string]io.Reader{
		"transactionName": strings.NewReader(td.Name),
		"amount":          strings.NewReader(strconv.FormatInt(td.Amount, 10)),
		"date":            strings.NewReader(strconv.FormatInt(td.Date, 10)),
	}
}

// NewCreateTransactionRequest creates a request to be used in tests to create a
// transaction, simulating data submitted from a form
func NewCreateTransactionRequest(values map[string]io.Reader) *http.Request {
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

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions", &b)
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

// NewGoogleCallbackRequest creates a request to be used in tests to call the
// Google callback route.
func NewGoogleCallbackRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/callback_google", nil)
	return req
}

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

// MakeCreateSharedTxnRequestPayload generates the payload to be given to NewCreateSharedTxnRequest
func MakeCreateSharedTxnRequestPayload(txn SharedTransaction) map[string]io.Reader {
	// make a comma separated list of participants
	participants := strings.Join(txn.Participants, ",")

	var unsettled string
	if txn.Unsettled {
		unsettled = "true"
	}

	return map[string]io.Reader{
		"shop":   strings.NewReader(txn.Shop),
		"amount": strings.NewReader(strconv.FormatInt(txn.Amount, 10)),
		// NOTE: currently, date will never be empty, change this?
		"date":         strings.NewReader(strconv.FormatInt(txn.Date, 10)),
		"participants": strings.NewReader(participants),
		"unsettled":    strings.NewReader(unsettled),
	}
}

// NewCreateSharedTxnRequest creates a request to be used in tests to create a
// shared transaction, simulating data submitted from a form
func NewCreateSharedTxnRequest(txn SharedTransaction) *http.Request {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	values := MakeCreateSharedTxnRequestPayload(txn)
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

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/trackers/%s/transactions", txn.Tracker), &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
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
