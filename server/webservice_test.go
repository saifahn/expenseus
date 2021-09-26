package expenseus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGetExpenseByID(t *testing.T) {
	store := StubExpenseStore{
		users: []User{},
		expenses: map[string]Expense{
			"1":    TestSeanExpense,
			"9281": TestTomomiExpense,
		},
	}
	webservice := &WebService{&store, &StubGoogleOauthConfig{}}

	t.Run("get an expense by id", func(t *testing.T) {
		request := NewGetExpenseRequest("1")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, got, TestSeanExpense)
	})

	t.Run("gets another expense by id", func(t *testing.T) {
		request := NewGetExpenseRequest("9281")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, got, TestTomomiExpense)
	})

	t.Run("returns 404 on non-existent expense", func(t *testing.T) {
		request := NewGetExpenseRequest("13371337")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetExpenseByUser(t *testing.T) {
	store := StubExpenseStore{
		users: []User{
			TestSeanUser,
			TestTomomiUser,
		},
		expenses: map[string]Expense{
			"1":    TestSeanExpense,
			"9281": TestTomomiExpense,
		},
	}
	webservice := NewWebService(&store, &StubGoogleOauthConfig{})

	t.Run("gets tomochi's expenses", func(t *testing.T) {
		request := NewGetExpensesByUsernameRequest(TestTomomiUser.Username)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpensesByUsername)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Len(t, got, 1)
		assert.Contains(t, got, TestTomomiExpense)
	})

	t.Run("gets saifahn's expenses", func(t *testing.T) {
		request := NewGetExpensesByUsernameRequest(TestSeanUser.Username)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpensesByUsername)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Len(t, got, 1)
		assert.Contains(t, got, TestSeanExpense)
	})
}

func TestCreateExpense(t *testing.T) {
	store := StubExpenseStore{
		users:    []User{},
		expenses: map[string]Expense{},
	}
	webservice := NewWebService(&store, &StubGoogleOauthConfig{})

	t.Run("creates a new expense on POST", func(t *testing.T) {
		request := NewCreateExpenseRequest("tomomi", "Test Expense")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.CreateExpense)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
		// this is technically actually testing implementation
		// I should just test that RecordExpense has been called correctly with the right thing, not the outcome
		assert.Len(t, store.expenses, 1)
	})
}

func TestGetAllExpenses(t *testing.T) {
	t.Run("gets all expenses with one expense", func(t *testing.T) {
		wantedExpenses := []Expense{
			TestTomomiExpense,
		}
		store := StubExpenseStore{
			users: []User{},
			expenses: map[string]Expense{
				"9281": TestTomomiExpense,
			},
		}
		webservice := NewWebService(&store, &StubGoogleOauthConfig{})

		request := NewGetAllExpensesRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetAllExpenses)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(wantedExpenses), len(got))
		assert.ElementsMatch(t, got, wantedExpenses)
	})

	t.Run("gets all expenses with more than one expense", func(t *testing.T) {
		wantedExpenses := []Expense{
			TestSeanExpense, TestTomomiExpense, TestTomomiExpense2,
		}
		store := StubExpenseStore{
			users: []User{},
			expenses: map[string]Expense{
				"1":     TestSeanExpense,
				"9281":  TestTomomiExpense,
				"14928": TestTomomiExpense2,
			},
		}
		webservice := NewWebService(&store, &StubGoogleOauthConfig{})

		request := NewGetAllExpensesRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetAllExpenses)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(wantedExpenses), len(got))
		assert.ElementsMatch(t, got, wantedExpenses)
	})

}

func TestCreateUser(t *testing.T) {
	store := StubExpenseStore{}
	webservice := NewWebService(&store, &StubGoogleOauthConfig{})

	user := TestSeanUser
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf(err.Error())
	}

	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("request could not be created, %v", err)
	}
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(webservice.CreateUser)
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusAccepted, response.Code)
	assert.Len(t, store.users, 1)
	assert.Contains(t, store.users, user)
}

func TestListUsers(t *testing.T) {
	store := StubExpenseStore{users: []User{TestSeanUser, TestTomomiUser}}
	webservice := NewWebService(&store, &StubGoogleOauthConfig{})

	request, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Fatalf("request could not be created, %v", err)
	}
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(webservice.ListUsers)
	handler.ServeHTTP(response, request)

	var got []User
	err = json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Len(t, got, len(store.users))
	assert.ElementsMatch(t, got, store.users)
}

func TestGoogleOauthCallback(t *testing.T) {
	t.Run("creates a user when user doesn't exist yet", func(t *testing.T) {
		store := StubExpenseStore{users: []User{}}
		oauth := StubGoogleOauthConfig{}
		webservice := NewWebService(&store, &oauth)

		request, err := http.NewRequest(http.MethodGet, "/callback_google", nil)
		if err != nil {
			t.Fatalf("request could not be created, %v", err)
		}
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.OauthCallback)
		handler.ServeHTTP(response, request)

		// expect a new user to be added to the store
		assert.Len(t, store.users, 1)
		// TODO: expect to be routed to update username page
	})
}

type StubGoogleOauthConfig struct{}

func (o *StubGoogleOauthConfig) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return ""
}

func (o *StubGoogleOauthConfig) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return nil, nil
}

func (o *StubGoogleOauthConfig) getInfoAndGenerateUser(state string, code string) (User, error) {
	return User{}, nil
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

func (s *StubExpenseStore) CreateUser(u User) error {
	s.users = append(s.users, u)
	return nil
}

func (s *StubExpenseStore) GetAllUsers() ([]User, error) {
	return s.users, nil
}
