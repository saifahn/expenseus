package expenseus_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/saifahn/expenseus"
	"github.com/saifahn/expenseus/ddb"
	"github.com/stretchr/testify/assert"
)

const usersTableName = "integtest-users-table"
const transactionsTableName = "integtest-transactions-table"

func setUpDB(d dynamodbiface.DynamoDBAPI) (expenseus.ExpenseStore, error) {
	err := ddb.CreateTestTable(d, usersTableName)
	if err != nil {
		return nil, err
	}
	err = ddb.CreateTestTable(d, transactionsTableName)
	if err != nil {
		return nil, err
	}

	return ddb.New(d, usersTableName, transactionsTableName), nil
}

func tearDownDB(d dynamodbiface.DynamoDBAPI) error {
	err := ddb.DeleteTable(d, usersTableName)
	if err != nil {
		return err
	}
	err = ddb.DeleteTable(d, transactionsTableName)
	if err != nil {
		return err
	}

	return nil
}

// setUpTestServer sets up a server with with the real routes and a test
// dynamodb instance, with stubs for the rest of the webservice
func setUpTestServer(t *testing.T) (http.Handler, func(t *testing.T)) {
	ddbLocal := ddb.NewDynamoDBLocalAPI()
	db, err := setUpDB(ddbLocal)
	if err != nil {
		t.Fatalf("could not set up the database: %v", err)
	}

	oauth := &expenseus.StubOauthConfig{}
	auth := &expenseus.StubSessionManager{}
	images := &expenseus.StubImageStore{}
	webservice := expenseus.NewWebService(db, oauth, auth, "", images)
	router := expenseus.InitRouter(webservice)

	return router, func(t *testing.T) {
		err := tearDownDB(ddbLocal)
		if err != nil {
			t.Fatalf("could not tear down the database: %v", err)
		}
	}
}

func createUser(t *testing.T, user expenseus.User, r http.Handler) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal the user JSON: %v", err)
	}
	response := httptest.NewRecorder()
	request := expenseus.NewCreateUserRequest(userJSON)
	request.AddCookie(&expenseus.ValidCookie)
	r.ServeHTTP(response, request)
	assert.Equal(t, http.StatusAccepted, response.Code)
}

func TestCreatingUsersAndRetrievingThem(t *testing.T) {
	t.Run("a valid cookie must be provided in order to create a user, GetSelf will read the cookie and attempt to get the user from the ID within, and a user can be retrieved by ID", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// TRY to create a user WITHOUT a valid cookie
		userJSON, err := json.Marshal(expenseus.TestSeanUser)
		if err != nil {
			t.Fatalf("failed to marshal the user JSON: %v", err)
		}
		response := httptest.NewRecorder()
		request := expenseus.NewCreateUserRequest(userJSON)
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusUnauthorized, response.Code)

		// use a VALID cookie
		createUser(t, expenseus.TestSeanUser, router)

		// TRY GetSelf with different ID in the cookie
		// should not work as the userID from the cookie does not exist
		response = httptest.NewRecorder()
		request = expenseus.NewGetSelfRequest()
		request.AddCookie(&http.Cookie{
			Name:  "session",
			Value: "not-real-id",
		})
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusNotFound, response.Code)

		// use a cookie with the SAME ID
		response = httptest.NewRecorder()
		request = expenseus.NewGetSelfRequest()
		request.AddCookie(&expenseus.ValidCookie)
		router.ServeHTTP(response, request)

		var userGot expenseus.User
		err = json.NewDecoder(response.Body).Decode(&userGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(expenseus.TestSeanUser, userGot)

		// GET the specifically created user from the db by ID
		response = httptest.NewRecorder()
		request = expenseus.NewGetUserRequest(expenseus.TestSeanUser.ID)
		request.AddCookie(&expenseus.ValidCookie)
		router.ServeHTTP(response, request)

		err = json.NewDecoder(response.Body).Decode(&userGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(expenseus.TestSeanUser, userGot)
	})

	t.Run("multiple users can be created and retrieved with a request to the GetAllUsers route", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create TWO users
		createUser(t, expenseus.TestSeanUser, router)
		createUser(t, expenseus.TestTomomiUser, router)

		// GET all users
		response := httptest.NewRecorder()
		request := expenseus.NewGetAllUsersRequest()
		request.AddCookie(&expenseus.ValidCookie)
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		// ensure that they contain the two users
		var usersGot []expenseus.User
		err := json.NewDecoder(response.Body).Decode(&usersGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Users: %v", response.Body, err)
		}
		assert.ElementsMatch(usersGot, []expenseus.User{expenseus.TestSeanUser, expenseus.TestTomomiUser})
	})
}

func TestCreatingExpensesAndRetrievingThem(t *testing.T) {
	var createTestExpense = func(t *testing.T, r http.Handler, ed expenseus.ExpenseDetails, userid string) {
		values := map[string]io.Reader{
			"expenseName": strings.NewReader(ed.Name),
		}
		request := expenseus.NewCreateExpenseRequest(values)
		request.AddCookie(&http.Cookie{Name: "session", Value: userid})
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)
		assert.Equal(t, http.StatusAccepted, response.Code)
	}

	t.Run("an expense can be added with a valid cookie and be retrieved as part of a GetAll request", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create user in the db
		createUser(t, expenseus.TestSeanUser, router)

		// create a transaction and store it
		wantedExpenseDetails := expenseus.TestSeanExpenseDetails
		createTestExpense(t, router, wantedExpenseDetails, wantedExpenseDetails.UserID)

		// try and get it
		request := expenseus.NewGetAllExpensesRequest()
		request.AddCookie(&http.Cookie{Name: "session", Value: wantedExpenseDetails.UserID})
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var expensesGot []expenseus.Expense
		err := json.NewDecoder(response.Body).Decode(&expensesGot)
		if err != nil {
			t.Logf("error parsing response from server %q into slice of Expenses: %v", response.Body, err)
		}

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(expensesGot, 1)
		assert.Equal(expensesGot[0].ExpenseDetails, wantedExpenseDetails)
	})

	t.Run("expenses can be retrieved by ID", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)
		createUser(t, expenseus.TestSeanUser, router)

		wantedExpenseDetails := expenseus.TestSeanExpenseDetails
		createTestExpense(t, router, wantedExpenseDetails, wantedExpenseDetails.UserID)

		request := expenseus.NewGetAllExpensesRequest()
		request.AddCookie(&expenseus.ValidCookie)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var expensesGot []expenseus.Expense
		err := json.NewDecoder(response.Body).Decode(&expensesGot)
		if err != nil {
			t.Logf("error parsing response from server %q into slice of Expenses: %v", response.Body, err)
		}

		// make sure the ID exists on the struct
		expenseID := expensesGot[0].ID
		assert.NotZero(expenseID)

		request = expenseus.NewGetExpenseRequest(expenseID)
		request.AddCookie(&expenseus.ValidCookie)
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		var got expenseus.Expense
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Errorf("error parsing response from server %q into Expense struct: %v", response.Body, err)
		}

		assert.Equal(wantedExpenseDetails, got.ExpenseDetails)
		assert.Equal(expensesGot[0], got)
	})

	// maybe just by user ID is better
	t.Run("expenses can be retrieved by username", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)
		createUser(t, expenseus.TestSeanUser, router)

		wantedExpenseDetails := expenseus.TestSeanExpenseDetails
		createTestExpense(t, router, wantedExpenseDetails, expenseus.TestSeanUser.ID)

		request := expenseus.NewGetExpensesByUsernameRequest(expenseus.TestSeanUser.Username)
		request.AddCookie(&expenseus.ValidCookie)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		var expensesGot []expenseus.Expense
		err := json.NewDecoder(response.Body).Decode(&expensesGot)
		if err != nil {
			t.Logf("error parsing response from server %q into slice of Expenses: %v", response.Body, err)
		}

		assert.Len(expensesGot, 1)
		assert.Equal(wantedExpenseDetails, expensesGot[0].ExpenseDetails)
	})
}
