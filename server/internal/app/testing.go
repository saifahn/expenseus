package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/nabeken/aws-go-dynamodb/table"
	"golang.org/x/oauth2"
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

// AddUserCookieContext adds a cookie and a user context to simulate a user
// being logged in.
func AddUserCookieAndContext(req *http.Request, id string) *http.Request {
	req.AddCookie(&http.Cookie{Name: "session", Value: id})
	ctx := context.WithValue(req.Context(), CtxKeyUserID, id)
	return req.WithContext(ctx)
}

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
func NewCreateTrackerRequest(t testing.TB, trackerDetails Tracker, userID string) *http.Request {
	trackerJSON, err := json.Marshal(trackerDetails)
	if err != nil {
		t.Fatalf(err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/trackers", bytes.NewBuffer(trackerJSON))
	// needs to pass in the userID of the user making the request to create the tracker
	ctx := context.WithValue(req.Context(), CtxKeyUserID, userID)
	return req.WithContext(ctx)
}

// NewGetTrackerByIDRequest creates a request to be used in tests to get a
// tracker by ID, with the ID in the request context.
func NewGetTrackerByIDRequest(t testing.TB, trackerID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/trackers/%s", trackerID), nil)
	ctx := context.WithValue(req.Context(), CtxKeyTrackerID, trackerID)
	return req.WithContext(ctx)
}

// NewGetTrackerByUserRequest creates a request to be used in tests to get a
// tracker by userID, with the userID in the request context.
func NewGetTrackerByUserRequest(t testing.TB, userID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/trackers/user/%s", userID), nil)
	ctx := context.WithValue(req.Context(), CtxKeyUserID, userID)
	return req.WithContext(ctx)
}

// NewGetTxnsByTrackerRequest creates a request to be used in tests to get a
// a list of transactions by trackerID, with the trackerID in the request context
func NewGetTxnsByTrackerRequest(t testing.TB, trackerID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/trackers/%s/transactions", trackerID), nil)
	ctx := context.WithValue(req.Context(), CtxKeyTrackerID, trackerID)
	return req.WithContext(ctx)
}

// MakeCreateSharedTxnRequestPayload generates the payload to be given to NewCreateSharedTxnRequest
func MakeCreateSharedTxnRequestPayload(txn SharedTransaction) map[string]io.Reader {
	// make a comma separated list of participants
	participants := strings.Join(txn.Participants, ",")

	return map[string]io.Reader{
		"shop":   strings.NewReader(txn.Shop),
		"amount": strings.NewReader(strconv.FormatInt(txn.Amount, 10)),
		// NOTE: currently, date will never be empty, change this?
		"date":         strings.NewReader(strconv.FormatInt(txn.Date, 10)),
		"participants": strings.NewReader(participants),
	}
}

// NewCreateSharedTxnRequest creates a request to be used in tests to create a
// shared transaction, simulating data submitted from a form
func NewCreateSharedTxnRequest(txn SharedTransaction, userID string) *http.Request {
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
	ctx := context.WithValue(req.Context(), CtxKeyUserID, userID)
	return req.WithContext(ctx)
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

// #region Sessions
type StubSessionManager struct {
	saveCalls   []string
	removeCalls int
}

var ValidCookie = http.Cookie{
	Name:  "session",
	Value: TestSeanUser.ID,
}

func (s *StubSessionManager) Validate(r *http.Request) bool {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == ValidCookie.Name {
			if len(cookie.Value) > 0 {
				return true
			}
		}
	}
	return false
}

func (s *StubSessionManager) Save(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)
	s.saveCalls = append(s.saveCalls, userID)
	http.SetCookie(rw, &ValidCookie)
}

func (s *StubSessionManager) GetUserID(r *http.Request) (string, error) {
	// get it from the cookie
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == ValidCookie.Name {
			return cookie.Value, nil
		}
	}
	return "", errors.New("no user ID was found")
}

func (s *StubSessionManager) Remove(rw http.ResponseWriter, r *http.Request) {
	s.removeCalls++
}

// #endregion Sessions

// #region OAuth
type StubOauthConfig struct {
	AuthCodeURLCalls []string
}

const oauthProviderMockURL = "oauth-provider-mock-url"

func (o *StubOauthConfig) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	o.AuthCodeURLCalls = append(o.AuthCodeURLCalls, oauthProviderMockURL)
	return oauthProviderMockURL
}

func (o *StubOauthConfig) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return nil, nil
}

func (o *StubOauthConfig) GetInfoAndGenerateUser(state string, code string) (User, error) {
	return TestSeanUser, nil
}

// #endregion OAuth

// #region Store
type StubTransactionStore struct {
	transactions                   map[string]Transaction
	users                          []User
	trackers                       []Tracker
	recordTransactionCalls         []TransactionDetails
	getTxnsByTrackerCalls          []string
	createSharedTxnCalls           []SharedTransaction
	getUnsettledTxnsByTrackerCalls []string
}

func (s *StubTransactionStore) GetTransaction(transactionID string) (Transaction, error) {
	transaction := s.transactions[transactionID]
	// check for empty Transaction
	if transaction == (Transaction{}) {
		return Transaction{}, table.ErrItemNotFound
	}
	return transaction, nil
}

func (s *StubTransactionStore) GetTransactionsByUser(id string) ([]Transaction, error) {
	var targetUser User
	for _, u := range s.users {
		if u.ID == id {
			targetUser = u
			break
		}
	}

	var transactions []Transaction
	for _, e := range s.transactions {
		// if the user id is the same as userid, then append
		if e.UserID == targetUser.ID {
			transactions = append(transactions, e)
		}
	}
	return transactions, nil
}

func (s *StubTransactionStore) CreateTransaction(ed TransactionDetails) error {
	testID := fmt.Sprintf("tid-%v", ed.Name)
	transaction := Transaction{
		TransactionDetails: ed,
		ID:                 testID,
	}
	s.transactions[testID] = transaction
	s.recordTransactionCalls = append(s.recordTransactionCalls, ed)
	return nil
}

func (s *StubTransactionStore) GetAllTransactions() ([]Transaction, error) {
	var transactions []Transaction
	for _, e := range s.transactions {
		transactions = append(transactions, e)
	}
	return transactions, nil
}

func (s *StubTransactionStore) GetUser(id string) (User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (s *StubTransactionStore) CreateUser(u User) error {
	s.users = append(s.users, u)
	return nil
}

func (s *StubTransactionStore) GetAllUsers() ([]User, error) {
	return s.users, nil
}

func (s *StubTransactionStore) CreateTracker(t Tracker) error {
	s.trackers = append(s.trackers, t)
	return nil
}

func (s *StubTransactionStore) GetTracker(id string) (Tracker, error) {
	for _, t := range s.trackers {
		if t.ID == id {
			return t, nil
		}
	}
	return Tracker{}, errors.New("tracker not found")
}

func (s *StubTransactionStore) GetTrackersByUser(userID string) ([]Tracker, error) {
	var trackers []Tracker
	for _, t := range s.trackers {
		for _, uid := range t.Users {
			if uid == userID {
				trackers = append(trackers, t)
			}
		}
	}
	return trackers, nil
}

func (s *StubTransactionStore) GetTxnsByTracker(trackerID string) ([]SharedTransaction, error) {
	s.getTxnsByTrackerCalls = append(s.getTxnsByTrackerCalls, trackerID)
	return nil, nil
}

func (s *StubTransactionStore) CreateSharedTxn(txn SharedTransaction) error {
	s.createSharedTxnCalls = append(s.createSharedTxnCalls, txn)
	return nil
}

func (s *StubTransactionStore) GetUnsettledTxnsByTracker(trackerID string) ([]SharedTransaction, error) {
	s.getUnsettledTxnsByTrackerCalls = append(s.getUnsettledTxnsByTrackerCalls, trackerID)
	return nil, nil
}

func (s *StubTransactionStore) SettleTxns(txns []SharedTransaction) error {
	return nil
}

// #endregion Store

// #region ImageStore
const testImageKey = "TEST_IMAGE_KEY"

type StubImageStore struct {
	uploadCalls                []string
	addImageToTransactionCalls []string
}

func (is *StubImageStore) Upload(file multipart.File, header multipart.FileHeader) (string, error) {
	is.uploadCalls = append(is.uploadCalls, "called")
	return testImageKey, nil
}

func (is *StubImageStore) Validate(file multipart.File) (bool, error) {
	return true, nil
}

func (is *StubImageStore) AddImageToTransaction(transaction Transaction) (Transaction, error) {
	is.addImageToTransactionCalls = append(is.addImageToTransactionCalls, "called")
	transaction.ImageURL = "test-image-url"
	return transaction, nil
}

// #endregion ImageStore

// #region InvalidImageStore
type StubInvalidImageStore struct {
	uploadCalls []string
}

func (is *StubInvalidImageStore) Upload(file multipart.File, header multipart.FileHeader) (string, error) {
	return "", errors.New("upload failed for some reason")
}

func (is *StubInvalidImageStore) Validate(file multipart.File) (bool, error) {
	return false, nil
}

func (is *StubInvalidImageStore) AddImageToTransaction(transaction Transaction) (Transaction, error) {
	return transaction, errors.New("image could not be added")
}

// #endregion InvalidImageStore
