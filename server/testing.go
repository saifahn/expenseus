package expenseus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

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
)

// NewGetExpenseRequest creates a request to be used in tests get an expense
// by id, adding the id to the request context.
func NewGetExpenseRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/%s", id), nil)
	ctx := context.WithValue(req.Context(), CtxKeyExpenseID, id)
	return req.WithContext(ctx)
}

// NewCreateExpenseRequest creates a request to be used in tests to create an
// expense that is associated with a user.
func NewCreateExpenseRequest(user, name string) *http.Request {
	values := ExpenseDetails{UserID: user, Name: name}
	jsonValue, _ := json.Marshal(values)
	req, _ := http.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(jsonValue))
	return req
}

// NewGetExpensesByUsernameRequest creates a request to be used in tests to get all
// expenses of a user, adding the user to the request context.
func NewGetExpensesByUsernameRequest(username string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/user/%s", username), nil)
	ctx := context.WithValue(req.Context(), CtxKeyUsername, username)
	return req.WithContext(ctx)
}

func NewGetAllExpensesRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/expenses", nil)
	return req
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got response body of %q, want %q", got, want)
	}
}

func AssertResponseStatus(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

type StubSessionManager struct {
	saveSessionCalls []string
}

var validCookie = http.Cookie{
	Name:  "session",
	Value: TestSeanUser.ID,
}

func (s *StubSessionManager) ValidateAuthorizedSession(r *http.Request) bool {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == validCookie.Name {
			if cookie.Value == validCookie.Value {
				return true
			}
		}
	}
	return false
}

func (s *StubSessionManager) SaveSession(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)
	s.saveSessionCalls = append(s.saveSessionCalls, userID)
	http.SetCookie(rw, &validCookie)
}

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

// stub store implementation
type StubExpenseStore struct {
	expenses map[string]Expense
	users    []User
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
