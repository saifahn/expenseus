package expenseus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExpenseByID(t *testing.T) {
	store := StubExpenseStore{
		users: []User{},
		expenses: map[string]Expense{
			"1":    TestSeanExpense,
			"9281": TestTomomiExpense,
		},
	}
	webservice := &WebService{&store, &StubOauthConfig{}, &StubSessionManager{}, ""}

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
	webservice := NewWebService(&store, &StubOauthConfig{}, &StubSessionManager{}, "")

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
	webservice := NewWebService(&store, &StubOauthConfig{}, &StubSessionManager{}, "")

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
		webservice := NewWebService(&store, &StubOauthConfig{}, &StubSessionManager{}, "")

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
		webservice := NewWebService(&store, &StubOauthConfig{}, &StubSessionManager{}, "")

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
	webservice := NewWebService(&store, &StubOauthConfig{}, &StubSessionManager{}, "")

	user := TestSeanUser
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf(err.Error())
	}

	request := NewCreateUserRequest(userJSON)
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(webservice.CreateUser)
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusAccepted, response.Code)
	assert.Len(t, store.users, 1)
	assert.Contains(t, store.users, user)
}

func TestListUsers(t *testing.T) {
	store := StubExpenseStore{users: []User{TestSeanUser, TestTomomiUser}}
	webservice := NewWebService(&store, &StubOauthConfig{}, &StubSessionManager{}, "")

	request := NewGetAllUsersRequest()
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(webservice.ListUsers)
	handler.ServeHTTP(response, request)

	var got []User
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Len(t, got, len(store.users))
	assert.ElementsMatch(t, got, store.users)
}

func TestOauthLogin(t *testing.T) {
	store := StubExpenseStore{}
	oauth := StubOauthConfig{}
	webservice := NewWebService(&store, &oauth, &StubSessionManager{}, "")

	request, err := http.NewRequest(http.MethodGet, "/api/v1/login_google", nil)
	if err != nil {
		t.Fatalf("request could not be created, %v", err)
	}
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(webservice.OauthLogin)
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusTemporaryRedirect, response.Code)
	// are these even good assertions to have?
	expectedURL := fmt.Sprintf("/api/v1/%s", oauthProviderMockURL)
	assert.Equal(t, expectedURL, response.Header().Get("Location"))
	// assert AuthCodeURL was called
	assert.Len(t, oauth.AuthCodeURLCalls, 1)
}

func TestOauthCallback(t *testing.T) {
	t.Run("creates a user when user doesn't exist yet and creates a new session with the user", func(t *testing.T) {
		store := StubExpenseStore{users: []User{}}
		oauth := StubOauthConfig{}
		sessions := StubSessionManager{}
		frontend := "http://a.test"
		webservice := NewWebService(&store, &oauth, &sessions, frontend)

		request := NewGoogleCallbackRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.OauthCallback)
		handler.ServeHTTP(response, request)

		// expect a new user to be added to the store, GetInfoAndGenerateUser has been stubbed to generate TestSeanUser
		expected := []User{TestSeanUser}
		assert.Len(t, store.users, 1)
		assert.ElementsMatch(t, expected, store.users)

		assert.Len(t, sessions.saveSessionCalls, 1)
		assert.Equal(t, sessions.saveSessionCalls[0], TestSeanUser.ID)

		// get routed to the base page for now
		url, err := response.Result().Location()
		if err != nil {
			t.Fatalf("url couldn't be found: %v", err)
		}
		assert.Equal(t, frontend, url.String())
	})

	t.Run("doesn't create a new user when the user already exists, and saves the session with the user in the context", func(t *testing.T) {
		store := StubExpenseStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		sessions := StubSessionManager{}
		frontend := "http://another.test"
		webservice := NewWebService(&store, &oauth, &sessions, frontend)

		request := NewGoogleCallbackRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(webservice.OauthCallback)
		handler.ServeHTTP(response, request)

		expected := []User{TestSeanUser}
		assert.Len(t, store.users, 1)
		assert.ElementsMatch(t, expected, store.users)

		assert.Len(t, sessions.saveSessionCalls, 1)
		// the callback will add a context of the appropriate user id
		assert.Equal(t, sessions.saveSessionCalls[0], TestSeanUser.ID)

		// expect to get routed to the main welcome page
		url, err := response.Result().Location()
		if err != nil {
			t.Fatalf("url couldn't be found: %v", err)
		}
		assert.Equal(t, frontend, url.String())
	})
}

func TestVerifyUser(t *testing.T) {
	t.Run("returns a 401 response when the user is not authorized", func(t *testing.T) {
		store := StubExpenseStore{}
		oauth := StubOauthConfig{}
		wb := NewWebService(&store, &oauth, &StubSessionManager{}, "")

		request := NewGetAllExpensesRequest()
		response := httptest.NewRecorder()

		handler := wb.VerifyUser(http.HandlerFunc(wb.GetAllExpenses))
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("returns a 200 response when the user is authorized, and passes the request to the appropriate route", func(t *testing.T) {
		store := StubExpenseStore{expenses: map[string]Expense{"1": TestSeanExpense}}
		oauth := StubOauthConfig{}
		wb := NewWebService(&store, &oauth, &StubSessionManager{}, "")

		request := NewGetAllExpensesRequest()
		// simulate a cookie session storage here
		request.AddCookie(&ValidCookie)
		response := httptest.NewRecorder()

		handler := wb.VerifyUser(http.HandlerFunc(wb.GetAllExpenses))
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.ElementsMatch(t, got, []Expense{TestSeanExpense})
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("returns a users details if the user exists", func(t *testing.T) {
		store := StubExpenseStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		wb := NewWebService(&store, &oauth, &StubSessionManager{}, "")

		request := NewGetUserRequest(TestSeanUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(wb.GetUser)
		handler.ServeHTTP(response, request)

		var got User
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, TestSeanUser, got)
	})
}

func TestGetSelf(t *testing.T) {
	t.Run("returns the user details from the stored session if the user exists", func(t *testing.T) {
		store := StubExpenseStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		wb := NewWebService(&store, &oauth, &StubSessionManager{}, "")

		request := NewGetSelfRequest()
		// add the user into the request cookie
		request.AddCookie(&ValidCookie)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(wb.GetSelf)
		handler.ServeHTTP(response, request)

		// decode the response
		var got User
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, TestSeanUser, got)
	})
}

func TestLogout(t *testing.T) {
	t.Run("session manager calls remove", func(t *testing.T) {
		store := StubExpenseStore{users: []User{TestSeanUser}}
		oauth := StubOauthConfig{}
		sessions := StubSessionManager{}
		frontend := "http://test.base"
		wb := NewWebService(&store, &oauth, &sessions, frontend)

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/logout", nil)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(wb.Logout)
		handler.ServeHTTP(response, request)

		assert.Equal(t, 1, sessions.removeCalls)
		assert.Equal(t, http.StatusTemporaryRedirect, response.Code)

		url, err := response.Result().Location()
		if err != nil {
			t.Fatalf("url couldn't be found: %v", err)
		}
		assert.Equal(t, frontend, url.String())
	})
}
