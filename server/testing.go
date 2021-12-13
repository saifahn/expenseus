package expenseus

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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

	TestSeanExpenseDetails = ExpenseDetails{
		Name:   "Expense 1",
		UserID: TestSeanUser.ID,
	}
	TestSeanExpense = Expense{
		ID:             "1",
		ExpenseDetails: TestSeanExpenseDetails,
	}

	TestTomomiExpenseDetails = ExpenseDetails{
		Name:   "Expense 2",
		UserID: TestTomomiUser.ID,
	}
	TestTomomiExpense = Expense{
		ID:             "2",
		ExpenseDetails: TestTomomiExpenseDetails,
	}

	TestTomomiExpense2Details = ExpenseDetails{
		Name:   "Expense 3",
		UserID: TestTomomiUser.ID,
	}
	TestTomomiExpense2 = Expense{
		ID:             "3",
		ExpenseDetails: TestTomomiExpense2Details,
	}

	TestExpenseWithImage = Expense{
		ID: "123",
		ExpenseDetails: ExpenseDetails{
			Name:     "ExpenseWithImage",
			UserID:   "an_ID",
			ImageKey: "test-image-key",
		},
	}
)

// NewGetExpenseRequest creates a request to be used in tests get an expense
// by ID, with ID in the request context.
func NewGetExpenseRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/expenses/%s", id), nil)
	ctx := context.WithValue(req.Context(), CtxKeyExpenseID, id)
	return req.WithContext(ctx)
}

// NewCreateExpenseRequest creates a request to be used in tests to create an
// expense that is associated with a user.
func NewCreateExpenseRequest(values map[string]io.Reader) *http.Request {
	// prepare FormData to submit
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

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/expenses", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// NewGetExpensesByUsernameRequest creates a request to be used in tests to get all
// expenses of a user, with the user in the request context.
func NewGetExpensesByUsernameRequest(username string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/expenses/user/%s", username), nil)
	ctx := context.WithValue(req.Context(), CtxKeyUsername, username)
	return req.WithContext(ctx)
}

// NewGetAllExpensesRequest creates a request to be used in tests to get all
// expenses.
func NewGetAllExpensesRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/expenses", nil)
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
type StubExpenseStore struct {
	expenses           map[string]Expense
	users              []User
	recordExpenseCalls []ExpenseDetails
}

func (s *StubExpenseStore) GetExpense(id string) (Expense, error) {
	expense := s.expenses[id]
	// check for empty Expense
	if expense == (Expense{}) {
		return Expense{}, errors.New("expense not found")
	}
	return expense, nil
}

func (s *StubExpenseStore) GetExpensesByUsername(username string) ([]Expense, error) {
	var targetUser User
	for _, u := range s.users {
		if u.Username == username {
			targetUser = u
			break
		}
	}

	var expenses []Expense
	for _, e := range s.expenses {
		// if the user id is the same as userid, then append
		if e.UserID == targetUser.ID {
			expenses = append(expenses, e)
		}
	}
	return expenses, nil
}

func (s *StubExpenseStore) RecordExpense(ed ExpenseDetails) error {
	testId := fmt.Sprintf("tid-%v", ed.Name)
	expense := Expense{
		ExpenseDetails: ed,
		ID:             testId,
	}
	s.expenses[testId] = expense
	s.recordExpenseCalls = append(s.recordExpenseCalls, ExpenseDetails{
		Name: ed.Name, UserID: ed.UserID, ImageKey: ed.ImageKey,
	})
	return nil
}

func (s *StubExpenseStore) GetAllExpenses() ([]Expense, error) {
	var expenses []Expense
	for _, e := range s.expenses {
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (s *StubExpenseStore) GetUser(id string) (User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (s *StubExpenseStore) CreateUser(u User) error {
	s.users = append(s.users, u)
	return nil
}

func (s *StubExpenseStore) GetAllUsers() ([]User, error) {
	return s.users, nil
}

// #endregion Store

// #region ImageStore
const testImageKey = "TEST_IMAGE_KEY"

type StubImageStore struct {
	uploadCalls            []string
	addImageToExpenseCalls []string
}

func (is *StubImageStore) Upload(file multipart.File, header multipart.FileHeader) (string, error) {
	is.uploadCalls = append(is.uploadCalls, "called")
	return testImageKey, nil
}

func (is *StubImageStore) Validate(file multipart.File) (bool, error) {
	return true, nil
}

func (is *StubImageStore) AddImageToExpense(expense Expense) (Expense, error) {
	is.addImageToExpenseCalls = append(is.addImageToExpenseCalls, "called")
	expense.ImageURL = "test-image-url"
	return expense, nil
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

func (is *StubInvalidImageStore) AddImageToExpense(expense Expense) (Expense, error) {
	return expense, errors.New("image could not be added")
}

// #endregion InvalidImageStore
