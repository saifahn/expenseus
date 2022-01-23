package expenseus_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nabeken/aws-go-dynamodb/table"
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
	uTbl := table.New(d, usersTableName)
	usersTable := ddb.NewUsersTable(uTbl)
	tTbl := table.New(d, transactionsTableName)
	transactionsTable := ddb.NewTransactionsTable(tTbl)

	return ddb.New(&usersTable, &transactionsTable), nil
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
	t.Run("an expense can be added with a valid cookie and be retrieved as part of a GetAll request", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create user in the db
		createUser(t, expenseus.TestSeanUser, router)

		// create a transaction and store it
		wantedExpenseDetails := expenseus.TestSeanExpenseDetails
		values := map[string]io.Reader{
			"expenseName": strings.NewReader(wantedExpenseDetails.Name),
		}
		request := expenseus.NewCreateExpenseRequest(values)
		request.AddCookie(&http.Cookie{Name: "session", Value: wantedExpenseDetails.UserID})
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusAccepted, response.Code)

		// try and get it
		request = expenseus.NewGetAllExpensesRequest()
		request.AddCookie(&http.Cookie{Name: "session", Value: wantedExpenseDetails.UserID})
		response = httptest.NewRecorder()
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
		values := map[string]io.Reader{
			"expenseName": strings.NewReader(wantedExpenseDetails.Name),
		}
		request := expenseus.NewCreateExpenseRequest(values)
		request.AddCookie(&http.Cookie{Name: "session", Value: wantedExpenseDetails.UserID})
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusAccepted, response.Code)

		request = expenseus.NewGetAllExpensesRequest()
		request.AddCookie(&expenseus.ValidCookie)
		response = httptest.NewRecorder()
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
			t.Errorf("error parsing response from server %q into struct of Expense: %v", response.Body, err)
		}

		assert.Equal(wantedExpenseDetails, got.ExpenseDetails)
		assert.Equal(expensesGot[0], got)
	})

	// maybe just by user ID is better
	t.Run("expenses can be retrieved by username", func(t *testing.T) {})

}

// func TestCreatingExpensesAndRetrievingThem(t *testing.T) {
// 	wantedExpenseDetails := []expenseus.ExpenseDetails{
// 		expenseus.TestSeanExpenseDetails,
// 		expenseus.TestTomomiExpenseDetails,
// 		expenseus.TestTomomiExpense2Details,
// 	}

// 	mr, err := miniredis.Run()
// 	if err != nil {
// 		t.Fatalf("error starting up miniredis, %v", err)
// 	}

// 	db := redis.New(mr.Addr())
// 	oauth := &expenseus.StubOauthConfig{}
// 	auth := &expenseus.StubSessionManager{}
// 	images := &expenseus.StubImageStore{}
// 	webservice := expenseus.NewWebService(db, oauth, auth, "", images)
// 	router := expenseus.InitRouter(webservice)

// 	// CREATE users in the db
// 	testUsers := []expenseus.User{expenseus.TestSeanUser, expenseus.TestTomomiUser}
// 	for _, u := range testUsers {
// 		userJSON, err := json.Marshal(u)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}

// 	response := httptest.NewRecorder()
// 	request := expenseus.NewCreateUserRequest(userJSON)
// 	request.AddCookie(&expenseus.ValidCookie)
// 	router.ServeHTTP(response, request)
// }

// 	// GET all users
// 	response := httptest.NewRecorder()
// 	request := expenseus.NewGetAllUsersRequest()
// 	request.AddCookie(&expenseus.ValidCookie)
// 	if err != nil {
// 		t.Fatalf("request could not be created, %v", err)
// 	}
// 	router.ServeHTTP(response, request)

// 	var usersGot []expenseus.User
// 	err = json.NewDecoder(response.Body).Decode(&usersGot)
// 	if err != nil {
// 		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
// 	}

// 	// ASSERT that the users received are correct
// 	assert.Len(t, usersGot, len(testUsers))
// 	assert.ElementsMatch(t, usersGot, testUsers)

// 	// CREATE expenses in the db
// 	for _, ed := range wantedExpenseDetails {
// 		values := map[string]io.Reader{
// 			"expenseName": strings.NewReader(ed.Name),
// 		}
// 		// TODO: should probably not use this method to be more like a real request
// 		request := expenseus.NewCreateExpenseRequest(values, ed.UserID)
// 		// VerifyUser now gets the UserID and passes it to the next handler, so
// 		// the proper UserID should be passed here.
// 		request.AddCookie(&http.Cookie{
// 			Name:  "session",
// 			Value: ed.UserID,
// 		})
// 		router.ServeHTTP(httptest.NewRecorder(), request)
// 	}

// 	// GET all expenses
// 	response = httptest.NewRecorder()
// 	request = expenseus.NewGetAllExpensesRequest()
// 	request.AddCookie(&expenseus.ValidCookie)
// 	router.ServeHTTP(response, request)

// 	var expensesGot []expenseus.Expense
// 	err = json.NewDecoder(response.Body).Decode(&expensesGot)
// 	if err != nil {
// 		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
// 	}

// 	assert.Equal(t, response.Code, http.StatusOK)
// 	assert.Equal(t, len(wantedExpenseDetails), len(expensesGot))
// 	// ASSERT there is an expense with the right details
// 	for _, expense := range expensesGot {
// 		assert.Contains(t, wantedExpenseDetails, expense.ExpenseDetails)
// 	}

// 	// GET one user's expenses
// 	response = httptest.NewRecorder()
// 	request = expenseus.NewGetExpensesByUsernameRequest(expenseus.TestTomomiUser.Username)
// 	request.AddCookie(&expenseus.ValidCookie)
// 	router.ServeHTTP(response, request)

// 	// reset expensesGot to an empty slice
// 	expensesGot = []expenseus.Expense{}
// 	err = json.NewDecoder(response.Body).Decode(&expensesGot)
// 	if err != nil {
// 		t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
// 	}

// 	assert.Len(t, expensesGot, 2)
// 	// ASSERT the expense details are the same. The ID is assigned separately, and we're currently not going to test that implementation
// 	var edsGot []expenseus.ExpenseDetails
// 	for _, e := range expensesGot {
// 		edsGot = append(edsGot, e.ExpenseDetails)
// 	}
// 	assert.Contains(t, edsGot, expenseus.TestTomomiExpense.ExpenseDetails)
// 	assert.Contains(t, edsGot, expenseus.TestTomomiExpense2.ExpenseDetails)
// }
