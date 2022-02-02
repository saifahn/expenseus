package app_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/saifahn/expenseus/internal/app"
	"github.com/saifahn/expenseus/internal/ddb"
	"github.com/saifahn/expenseus/internal/router"
	"github.com/stretchr/testify/assert"
)

const usersTableName = "integtest-users-table"
const transactionsTableName = "integtest-transactions-table"

func setUpDB(d dynamodbiface.DynamoDBAPI) (app.Store, error) {
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
// dynamodb instance, with stubs for the rest of the app
func setUpTestServer(t *testing.T) (http.Handler, func(t *testing.T)) {
	ddbLocal := ddb.NewDynamoDBLocalAPI()
	db, err := setUpDB(ddbLocal)
	if err != nil {
		t.Fatalf("could not set up the database: %v", err)
	}

	oauth := &app.StubOauthConfig{}
	auth := &app.StubSessionManager{}
	images := &app.StubImageStore{}
	a := app.New(db, oauth, auth, "", images)
	r := router.Init(a)

	return r, func(t *testing.T) {
		err := tearDownDB(ddbLocal)
		if err != nil {
			t.Fatalf("could not tear down the database: %v", err)
		}
	}
}

func createUser(t *testing.T, user app.User, r http.Handler) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal the user JSON: %v", err)
	}
	response := httptest.NewRecorder()
	request := app.NewCreateUserRequest(userJSON)
	request.AddCookie(&app.ValidCookie)
	r.ServeHTTP(response, request)
	assert.Equal(t, http.StatusAccepted, response.Code)
}

func TestCreatingUsersAndRetrievingThem(t *testing.T) {
	t.Run("a valid cookie must be provided in order to create a user, GetSelf will read the cookie and attempt to get the user from the ID within, and a user can be retrieved by ID", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// TRY to create a user WITHOUT a valid cookie
		userJSON, err := json.Marshal(app.TestSeanUser)
		if err != nil {
			t.Fatalf("failed to marshal the user JSON: %v", err)
		}
		response := httptest.NewRecorder()
		request := app.NewCreateUserRequest(userJSON)
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusUnauthorized, response.Code)

		// use a VALID cookie
		createUser(t, app.TestSeanUser, router)

		// TRY GetSelf with different ID in the cookie
		// should not work as the userID from the cookie does not exist
		response = httptest.NewRecorder()
		request = app.NewGetSelfRequest()
		request.AddCookie(&http.Cookie{
			Name:  "session",
			Value: "not-real-id",
		})
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusNotFound, response.Code)

		// use a cookie with the SAME ID
		response = httptest.NewRecorder()
		request = app.NewGetSelfRequest()
		request.AddCookie(&app.ValidCookie)
		router.ServeHTTP(response, request)

		var userGot app.User
		err = json.NewDecoder(response.Body).Decode(&userGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(app.TestSeanUser, userGot)

		// GET the specifically created user from the db by ID
		response = httptest.NewRecorder()
		request = app.NewGetUserRequest(app.TestSeanUser.ID)
		request.AddCookie(&app.ValidCookie)
		router.ServeHTTP(response, request)

		err = json.NewDecoder(response.Body).Decode(&userGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(app.TestSeanUser, userGot)
	})

	t.Run("multiple users can be created and retrieved with a request to the GetAllUsers route", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create TWO users
		createUser(t, app.TestSeanUser, router)
		createUser(t, app.TestTomomiUser, router)

		// GET all users
		response := httptest.NewRecorder()
		request := app.NewGetAllUsersRequest()
		request.AddCookie(&app.ValidCookie)
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		// ensure that they contain the two users
		var usersGot []app.User
		err := json.NewDecoder(response.Body).Decode(&usersGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Users: %v", response.Body, err)
		}
		assert.ElementsMatch(usersGot, []app.User{app.TestSeanUser, app.TestTomomiUser})
	})
}

func TestCreatingTransactionsAndRetrievingThem(t *testing.T) {
	var createTestTransaction = func(t *testing.T, r http.Handler, ed app.TransactionDetails, userid string) {
		payload := app.MakeCreateTransactionRequestPayload(ed)
		request := app.NewCreateTransactionRequest(payload)
		request.AddCookie(&http.Cookie{Name: "session", Value: userid})
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)
		assert.Equal(t, http.StatusAccepted, response.Code)
	}

	t.Run("an transaction can be added with a valid cookie and be retrieved as part of a GetAll request", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create user in the db
		createUser(t, app.TestSeanUser, router)

		// create a transaction and store it
		wantedTransactionDetails := app.TestSeanTransactionDetails
		createTestTransaction(t, router, wantedTransactionDetails, wantedTransactionDetails.UserID)

		// try and get it
		request := app.NewGetAllTransactionsRequest()
		request.AddCookie(&http.Cookie{Name: "session", Value: wantedTransactionDetails.UserID})
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var transactionsGot []app.Transaction
		err := json.NewDecoder(response.Body).Decode(&transactionsGot)
		if err != nil {
			t.Logf("error parsing response from server %q into slice of Transactions: %v", response.Body, err)
		}

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(transactionsGot, 1)
		assert.Equal(transactionsGot[0].TransactionDetails, wantedTransactionDetails)
	})

	t.Run("transactions can be retrieved by ID", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)
		createUser(t, app.TestSeanUser, router)

		wantedTransactionDetails := app.TestSeanTransactionDetails
		createTestTransaction(t, router, wantedTransactionDetails, wantedTransactionDetails.UserID)

		request := app.NewGetAllTransactionsRequest()
		request.AddCookie(&app.ValidCookie)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var transactionsGot []app.Transaction
		err := json.NewDecoder(response.Body).Decode(&transactionsGot)
		if err != nil {
			t.Logf("error parsing response from server %q into slice of Transactions: %v", response.Body, err)
		}

		// make sure the ID exists on the struct
		transactionID := transactionsGot[0].ID
		assert.NotZero(transactionID)

		request = app.NewGetTransactionRequest(transactionID)
		request.AddCookie(&app.ValidCookie)
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		var got app.Transaction
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Errorf("error parsing response from server %q into Transaction struct: %v", response.Body, err)
		}

		assert.Equal(wantedTransactionDetails, got.TransactionDetails)
		assert.Equal(transactionsGot[0], got)
	})

	// maybe just by user ID is better
	t.Run("transactions can be retrieved by username", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)
		createUser(t, app.TestSeanUser, router)

		wantedTransactionDetails := app.TestSeanTransactionDetails
		createTestTransaction(t, router, wantedTransactionDetails, app.TestSeanUser.ID)

		request := app.NewGetTransactionsByUsernameRequest(app.TestSeanUser.Username)
		request.AddCookie(&app.ValidCookie)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		var transactionsGot []app.Transaction
		err := json.NewDecoder(response.Body).Decode(&transactionsGot)
		if err != nil {
			t.Logf("error parsing response from server %q into slice of Transactions: %v", response.Body, err)
		}

		assert.Len(transactionsGot, 1)
		assert.Equal(wantedTransactionDetails, transactionsGot[0].TransactionDetails)
	})
}
