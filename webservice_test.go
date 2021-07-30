package expenseus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var seanUser = User{
	Username: "saifahn",
	Name:     "Sean Li",
	ID:       "sean_id",
}

var tomomiUser = User{
	Username: "tomochi",
	Name:     "Tomomi Kinoshita",
	ID:       "tomomi_id",
}

var testSeanExpense = Expense{
	ID:     "1",
	Name:   "Expense 1",
	UserID: seanUser.ID,
}

var testTomomiExpense = Expense{
	ID:     "9281",
	Name:   "Expense 9281",
	UserID: tomomiUser.ID,
}

var testTomomiExpense2 = Expense{
	ID:     "14928",
	Name:   "Expense 14928",
	UserID: tomomiUser.ID,
}

func TestGetExpenseByID(t *testing.T) {
	store := StubExpenseStore{
		users: []User{},
		expenses: map[string]Expense{
			"1":    testSeanExpense,
			"9281": testTomomiExpense,
		},
	}
	webservice := &WebService{&store}

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
		assert.Equal(t, got, testSeanExpense)
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
		assert.Equal(t, got, testTomomiExpense)
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
			seanUser,
			tomomiUser,
		},
		expenses: map[string]Expense{
			"1":    testSeanExpense,
			"9281": testTomomiExpense,
		},
	}
	webservice := NewWebService(&store)

	t.Run("gets tomochi's expenses", func(t *testing.T) {
		request := NewGetExpensesByUsernameRequest(tomomiUser.Username)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpensesByUser)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Len(t, got, 1)
		assert.Contains(t, got, testTomomiExpense)
	})

	t.Run("gets saifahn's expenses", func(t *testing.T) {
		request := NewGetExpensesByUsernameRequest(seanUser.Username)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.GetExpensesByUser)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Len(t, got, 1)
		assert.Contains(t, got, testSeanExpense)
	})
}

func TestCreateExpense(t *testing.T) {
	store := StubExpenseStore{
		users:    []User{},
		expenses: map[string]Expense{},
	}
	webservice := NewWebService(&store)

	t.Run("creates a new expense on POST", func(t *testing.T) {
		request := NewCreateExpenseRequest("tomomi", "Test Expense")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.CreateExpense)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
		assert.Len(t, store.expenses, 1)
	})
}

func TestGetAllExpenses(t *testing.T) {
	t.Run("gets all expenses with one expense", func(t *testing.T) {
		wantedExpenses := []Expense{
			testTomomiExpense,
		}
		store := StubExpenseStore{
			users: []User{},
			expenses: map[string]Expense{
				"9281": testTomomiExpense,
			},
		}
		webservice := NewWebService(&store)

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
			testSeanExpense, testTomomiExpense, testTomomiExpense2,
		}
		store := StubExpenseStore{
			users: []User{},
			expenses: map[string]Expense{
				"1":     testSeanExpense,
				"9281":  testTomomiExpense,
				"14928": testTomomiExpense2,
			},
		}
		webservice := NewWebService(&store)

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
	webservice := NewWebService(&store)

	user := User{ID: "saifahn", Name: "Sean Li"}
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

func (s *StubExpenseStore) GetExpensesByUser(username string) ([]Expense, error) {
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

func (s *StubExpenseStore) RecordExpense(e Expense) error {
	testId := fmt.Sprintf("tid-%v", e.Name)
	s.expenses[testId] = e
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
