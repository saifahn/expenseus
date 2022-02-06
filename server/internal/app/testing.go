package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

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
)

// addUserCookieContext adds a cookie and a user context to simulate a user
// being logged in.
func addUserCookieAndContext(req *http.Request, id string) *http.Request {
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
// CreateTransactionRequest
func MakeCreateTransactionRequestPayload(td TransactionDetails) map[string]io.Reader {
	return map[string]io.Reader{
		"transactionName": strings.NewReader(td.Name),
		"amount":          strings.NewReader(strconv.FormatInt(td.Amount, 10)),
		"date":            strings.NewReader(strconv.FormatInt(td.Date, 10)),
	}
}

// NewCreateTransactionRequest creates a request to be used in tests to create an
// transaction
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

// NewGetTransactionsByUsernameRequest creates a request to be used in tests to get all
// transactions of a user, with the user in the request context.
func NewGetTransactionsByUsernameRequest(username string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/transactions/user/%s", username), nil)
	ctx := context.WithValue(req.Context(), CtxKeyUsername, username)
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
	transactions           map[string]Transaction
	users                  []User
	recordTransactionCalls []TransactionDetails
}

func (s *StubTransactionStore) GetTransaction(id string) (Transaction, error) {
	transaction := s.transactions[id]
	// check for empty Transaction
	if transaction == (Transaction{}) {
		return Transaction{}, errors.New("transaction not found")
	}
	return transaction, nil
}

func (s *StubTransactionStore) GetTransactionsByUsername(username string) ([]Transaction, error) {
	var targetUser User
	for _, u := range s.users {
		if u.Username == username {
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
	testId := fmt.Sprintf("tid-%v", ed.Name)
	transaction := Transaction{
		TransactionDetails: ed,
		ID:                 testId,
	}
	s.transactions[testId] = transaction
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
