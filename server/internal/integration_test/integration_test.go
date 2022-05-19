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

const testTableName = "expenseus-integ-test"

func setUpDB(d dynamodbiface.DynamoDBAPI) (app.Store, error) {
	err := ddb.CreateTestTable(d, testTableName)
	if err != nil {
		return nil, err
	}

	return ddb.New(d, testTableName), nil
}

func tearDownDB(d dynamodbiface.DynamoDBAPI) error {
	err := ddb.DeleteTable(d, testTableName)
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
		assert.Equal(wantedTransactionDetails, transactionsGot[0].TransactionDetails)
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

	t.Run("transactions can be retrieved by user ID", func(t *testing.T) {
		router, tearDownDB := setUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)
		createUser(t, app.TestSeanUser, router)

		wantedTransactionDetails := app.TestSeanTransactionDetails
		createTestTransaction(t, router, wantedTransactionDetails, app.TestSeanUser.ID)

		request := app.NewGetTransactionsByUserRequest(app.TestSeanUser.ID)
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

func createTracker(t testing.TB, tracker app.Tracker, r http.Handler) {
	response := httptest.NewRecorder()
	request := app.NewCreateTrackerRequest(t, tracker)
	validCookie := http.Cookie{
		Name: "session",
		// the cookie has to contain information from one of the users in the tracker
		Value: tracker.Users[0],
	}
	request.AddCookie(&validCookie)
	r.ServeHTTP(response, request)
}

func TestCreatingTrackers(t *testing.T) {
	tests := map[string]struct {
		tracker      app.Tracker
		cookie       http.Cookie
		expectedCode int
	}{
		"without a valid cookie": {
			tracker:      app.TestTracker,
			cookie:       http.Cookie{Name: "invalid"},
			expectedCode: http.StatusUnauthorized,
		},
		"session user is not involved in tracker": {
			tracker:      app.TestTracker,
			cookie:       http.Cookie{Name: "session", Value: "not-in-tracker-user"},
			expectedCode: http.StatusForbidden,
		},
		"session is involved in tracker": {
			tracker:      app.TestTracker,
			cookie:       http.Cookie{Name: "session", Value: app.TestSeanTransaction.UserID},
			expectedCode: http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := setUpTestServer(t)
			defer tearDownDB(t)
			assert := assert.New(t)

			request := app.NewCreateTrackerRequest(t, tc.tracker)
			request.AddCookie(&tc.cookie)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.expectedCode, response.Code)
		})
	}
}

func TestGetTracker(t *testing.T) {
	router, tearDownDB := setUpTestServer(t)
	defer tearDownDB(t)
	assert := assert.New(t)

	tests := map[string]struct {
		trackerID    string
		cookie       http.Cookie
		expectedCode int
	}{
		"without a valid cookie": {
			trackerID:    "invalid",
			cookie:       http.Cookie{Name: "invalid"},
			expectedCode: http.StatusUnauthorized,
		},
		"with a non-existent tracker ID": {
			trackerID:    "non-existent-tracker-id",
			cookie:       app.ValidCookie,
			expectedCode: http.StatusNotFound,
		},
		// NOTE: we can't actually do this here because we don't know the ID
		// "with a tracker ID of an existing tracker": {
		// 	trackerID:    app.TestTracker.ID,
		// 	cookie:       app.ValidCookie,
		// 	expectedCode: http.StatusAccepted,
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request := app.NewGetTrackerByIDRequest(t, tc.trackerID)
			request.AddCookie(&tc.cookie)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.expectedCode, response.Code)
		})
	}
}

func TestGetTrackersByUser(t *testing.T) {
	router, tearDownDB := setUpTestServer(t)
	defer tearDownDB(t)
	assert := assert.New(t)
	createTracker(t, app.TestTracker, router)
	testTrackerWithTwoUsers := app.Tracker{
		Name:  "tracker for two",
		Users: []string{app.TestSeanUser.ID, app.TestTomomiUser.ID},
	}
	createTracker(t, testTrackerWithTwoUsers, router)

	tests := map[string]struct {
		user         string
		cookie       http.Cookie
		wantCode     int
		wantTrackers []app.Tracker
	}{
		"without a valid cookie": {
			user:         "invalid",
			cookie:       http.Cookie{Name: "invalid"},
			wantCode:     http.StatusUnauthorized,
			wantTrackers: nil,
		},
		"with a user in no trackers": {
			user:         "notInAnyTrackers",
			cookie:       app.ValidCookie,
			wantCode:     http.StatusOK,
			wantTrackers: nil,
		},
		"with a user in a tracker": {
			user:         app.TestTomomiUser.ID,
			cookie:       app.ValidCookie,
			wantCode:     http.StatusOK,
			wantTrackers: []app.Tracker{testTrackerWithTwoUsers},
		},
		"with a user in two trackers": {
			user:         app.TestSeanUser.ID,
			cookie:       app.ValidCookie,
			wantCode:     http.StatusOK,
			wantTrackers: []app.Tracker{app.TestTracker, testTrackerWithTwoUsers},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request := app.NewGetTrackerByUserRequest(t, tc.user)
			request.AddCookie(&tc.cookie)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)

			var gotTrackers []app.Tracker
			err := json.NewDecoder(response.Body).Decode(&gotTrackers)
			if err != nil {
				t.Logf("error parsing response from server %q into slice of Trackers: %v", response.Body, err)
			}
			assert.Len(gotTrackers, len(tc.wantTrackers))

			// check without ID because it will be set as a UUID in the real db
			var wantTrackersNoID, gotTrackersNoID []app.Tracker
			for _, wt := range tc.wantTrackers {
				wantTrackersNoID = append(wantTrackersNoID, app.Tracker{
					Name:  wt.Name,
					Users: wt.Users,
				})
			}
			for _, gt := range gotTrackers {
				gotTrackersNoID = append(gotTrackersNoID, app.Tracker{
					Name:  gt.Name,
					Users: gt.Users,
				})
			}
			assert.ElementsMatch(wantTrackersNoID, gotTrackersNoID)
		})
	}
}

func TestCreateSharedTxn(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		transaction app.SharedTransaction
		cookie      http.Cookie
		wantCode    int
	}{
		"without a valid cookie": {
			transaction: app.SharedTransaction{},
			cookie:      http.Cookie{Name: "invalid"},
			wantCode:    http.StatusUnauthorized,
		},
		"with a valid cookie but an invalid transaction": {
			transaction: app.SharedTransaction{},
			cookie:      app.ValidCookie,
			wantCode:    http.StatusBadRequest,
		},
		"with a valid cookie, but the user is not in the participants field": {
			transaction: app.SharedTransaction{
				Participants: []string{"user-01", "user-02"},
				Shop:         "test-shop",
				Amount:       123,
				Date:         123456,
			},
			cookie:   http.Cookie{Name: "session", Value: "not-in-participants"},
			wantCode: http.StatusForbidden,
		},
		"with a valid cookie and a valid transaction": {
			transaction: app.SharedTransaction{
				Participants: []string{"user-01", "user-02"},
				Shop:         "test-shop",
				Amount:       123,
				Date:         123456,
			},
			cookie:   http.Cookie{Name: "session", Value: "user-01"},
			wantCode: http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := setUpTestServer(t)
			defer tearDownDB(t)
			request := app.NewCreateSharedTxnRequest(tc.transaction)
			request.AddCookie(&tc.cookie)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}
