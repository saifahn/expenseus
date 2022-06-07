package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifahn/expenseus/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestCreatingUsersAndRetrievingThem(t *testing.T) {
	t.Run("a valid cookie must be provided in order to create a user, GetSelf will read the cookie and attempt to get the user from the ID within, and a user can be retrieved by ID", func(t *testing.T) {
		router, tearDownDB := SetUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// TRY to create a user WITHOUT a valid cookie
		userJSON, err := json.Marshal(TestSeanUser)
		if err != nil {
			t.Fatalf("failed to marshal the user JSON: %v", err)
		}
		response := httptest.NewRecorder()
		request := app.NewCreateUserRequest(userJSON)
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusUnauthorized, response.Code)

		// use a VALID cookie
		CreateUser(t, TestSeanUser, router)

		// TRY GetSelf with different ID in the cookie
		// should not work as the userID from the cookie does not exist
		response = httptest.NewRecorder()
		request = app.NewGetSelfRequest()
		request.AddCookie(CreateCookie("not-real-id"))

		router.ServeHTTP(response, request)
		assert.Equal(http.StatusNotFound, response.Code)

		// use a cookie with the SAME ID
		response = httptest.NewRecorder()
		request = app.NewGetSelfRequest()
		request.AddCookie(CreateCookie(TestSeanUser.ID))
		router.ServeHTTP(response, request)

		var userGot app.User
		err = json.NewDecoder(response.Body).Decode(&userGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(TestSeanUser, userGot)

		// GET the specifically created user from the db by ID
		response = httptest.NewRecorder()
		request = app.NewGetUserRequest(TestSeanUser.ID)
		request.AddCookie(CreateCookie(TestSeanUser.ID))
		router.ServeHTTP(response, request)

		err = json.NewDecoder(response.Body).Decode(&userGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into User struct, '%v'", response.Body, err)
		}
		assert.Equal(http.StatusOK, response.Code)
		assert.Equal(TestSeanUser, userGot)
	})

	t.Run("multiple users can be created and retrieved with a request to the GetAllUsers route", func(t *testing.T) {
		router, tearDownDB := SetUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create TWO users
		CreateUser(t, TestSeanUser, router)
		CreateUser(t, TestTomomiUser, router)

		// GET all users
		response := httptest.NewRecorder()
		request := app.NewGetAllUsersRequest()
		request.AddCookie(CreateCookie(TestSeanUser.ID))
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		// ensure that they contain the two users
		var usersGot []app.User
		err := json.NewDecoder(response.Body).Decode(&usersGot)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Users: %v", response.Body, err)
		}
		assert.ElementsMatch(usersGot, []app.User{TestSeanUser, TestTomomiUser})
	})
}

func createTestTransaction(t *testing.T, r http.Handler, td app.Transaction, userid string) {
	payload := app.MakeCreateTransactionRequestPayload(td)
	request := app.NewCreateTransactionRequest(payload)
	request.AddCookie(CreateCookie(userid))
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	assert.Equal(t, http.StatusAccepted, response.Code)
}

func TestCreatingTransactionsAndRetrievingThem(t *testing.T) {
	t.Run("an transaction can be added with a valid cookie and be retrieved as part of a GetAll request", func(t *testing.T) {
		router, tearDownDB := SetUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create user in the db
		CreateUser(t, TestSeanUser, router)

		// create a transaction and store it
		wantTxnDetails := TestSeanTxnDetails
		createTestTransaction(t, router, wantTxnDetails, wantTxnDetails.UserID)

		// try and get it
		request := app.NewGetAllTransactionsRequest()
		request.AddCookie(CreateCookie(wantTxnDetails.UserID))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var transactionsGot []app.Transaction
		err := json.NewDecoder(response.Body).Decode(&transactionsGot)
		if err != nil {
			t.Logf("error parsing response from server %q into slice of Transactions: %v", response.Body, err)
		}

		assert.Equal(http.StatusOK, response.Code)
		assert.Len(transactionsGot, 1)

		AssertEqualTxnDetails(t, wantTxnDetails, transactionsGot[0])
	})

	t.Run("transactions can be retrieved by ID", func(t *testing.T) {
		router, tearDownDB := SetUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)
		CreateUser(t, TestSeanUser, router)

		wantTxnDetails := TestSeanTxnDetails
		createTestTransaction(t, router, wantTxnDetails, wantTxnDetails.UserID)

		request := app.NewGetAllTransactionsRequest()
		request.AddCookie(CreateCookie(TestSeanUser.ID))
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
		request.AddCookie(CreateCookie(TestSeanUser.ID))
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		var got app.Transaction
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Errorf("error parsing response from server %q into Transaction struct: %v", response.Body, err)
		}

		AssertEqualTxnDetails(t, wantTxnDetails, got)
		assert.Equal(transactionsGot[0], got)
	})

	t.Run("transactions can be retrieved by user ID", func(t *testing.T) {
		router, tearDownDB := SetUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)
		CreateUser(t, TestSeanUser, router)

		wantTxnDetails := TestSeanTxnDetails
		createTestTransaction(t, router, wantTxnDetails, TestSeanUser.ID)

		request := app.NewGetTransactionsByUserRequest(TestSeanUser.ID)
		request.AddCookie(CreateCookie(TestSeanUser.ID))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(http.StatusOK, response.Code)

		var transactionsGot []app.Transaction
		err := json.NewDecoder(response.Body).Decode(&transactionsGot)
		if err != nil {
			t.Logf("error parsing response from server %q into slice of Transactions: %v", response.Body, err)
		}

		assert.Len(transactionsGot, 1)
		AssertEqualTxnDetails(t, wantTxnDetails, transactionsGot[0])
	})
}

func TestDeletingTransactions(t *testing.T) {
	tests := map[string]struct {
		td       app.Transaction
		user     string
		wantCode int
	}{
		"with a valid cookie, the transaction is deleted": {
			td:       TestSeanTxnDetails,
			user:     TestSeanUser.ID,
			wantCode: http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)
			assert := assert.New(t)

			CreateUser(t, TestSeanUser, router)
			createTestTransaction(t, router, tc.td, tc.user)

			request := app.NewGetTransactionsByUserRequest(tc.user)
			request.AddCookie(CreateCookie(tc.user))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			var got []app.Transaction
			err := json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Errorf("error parsing response from server %q into slice of Transactions: %v", response.Body, err)
			}
			assert.Len(got, 1)

			request = app.NewDeleteTransactionRequest(got[0].ID)
			request.AddCookie(CreateCookie(tc.user))
			response = httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(tc.wantCode, response.Code)

			// get the transactions again - this time it should be empty
			request = app.NewGetTransactionsByUserRequest(tc.user)
			request.AddCookie(CreateCookie(tc.user))
			response = httptest.NewRecorder()
			router.ServeHTTP(response, request)
			err = json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Errorf("error parsing response from server %q into slice of Transactions: %v", response.Body, err)
			}
			assert.Len(got, 0)
		})
	}
}

func createTracker(t testing.TB, tracker app.Tracker, r http.Handler) {
	response := httptest.NewRecorder()
	request := app.NewCreateTrackerRequest(t, tracker)
	// the cookie has to contain information from one of the users in the tracker
	request.AddCookie(CreateCookie(tracker.Users[0]))
	r.ServeHTTP(response, request)
}

func TestCreatingTrackers(t *testing.T) {
	tests := map[string]struct {
		tracker      app.Tracker
		cookie       http.Cookie
		expectedCode int
	}{
		"without a valid cookie": {
			tracker:      TestTracker,
			cookie:       http.Cookie{Name: "invalid"},
			expectedCode: http.StatusUnauthorized,
		},
		"session user is not involved in tracker": {
			tracker:      TestTracker,
			cookie:       *CreateCookie("not-in-tracker-user"),
			expectedCode: http.StatusForbidden,
		},
		"session is involved in tracker": {
			tracker:      TestTracker,
			cookie:       *CreateCookie(TestSeanTxnDetails.UserID),
			expectedCode: http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
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
	router, tearDownDB := SetUpTestServer(t)
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
			cookie:       *CreateCookie(TestSeanUser.ID),
			expectedCode: http.StatusNotFound,
		},
		// NOTE: we can't actually do this here because we don't know the ID
		// "with a tracker ID of an existing tracker": {
		// 	trackerID:    TestTracker.ID,
		// 	cookie:       app.ValidCookie,
		// 	expectedCode: http.StatusAccepted,
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request := app.NewGetTrackerByIDRequest(tc.trackerID)
			request.AddCookie(&tc.cookie)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.expectedCode, response.Code)
		})
	}
}

func TestGetTrackersByUser(t *testing.T) {
	router, tearDownDB := SetUpTestServer(t)
	defer tearDownDB(t)
	assert := assert.New(t)
	createTracker(t, TestTracker, router)
	testTrackerWithTwoUsers := app.Tracker{
		Name:  "tracker for two",
		Users: []string{TestSeanUser.ID, TestTomomiUser.ID},
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
			cookie:       *CreateCookie(TestSeanUser.ID),
			wantCode:     http.StatusOK,
			wantTrackers: nil,
		},
		"with a user in a tracker": {
			user:         TestTomomiUser.ID,
			cookie:       *CreateCookie(TestSeanUser.ID),
			wantCode:     http.StatusOK,
			wantTrackers: []app.Tracker{testTrackerWithTwoUsers},
		},
		"with a user in two trackers": {
			user:         TestSeanUser.ID,
			cookie:       *CreateCookie(TestSeanUser.ID),
			wantCode:     http.StatusOK,
			wantTrackers: []app.Tracker{TestTracker, testTrackerWithTwoUsers},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request := app.NewGetTrackerByUserRequest(tc.user)
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
			cookie:      *CreateCookie(TestSeanUser.ID),
			wantCode:    http.StatusBadRequest,
		},
		"with a valid cookie, but the user is not in the participants field": {
			transaction: app.SharedTransaction{
				Participants: []string{"user-01", "user-02"},
				Shop:         "test-shop",
				Amount:       123,
				Date:         123456,
			},
			cookie:   *CreateCookie("not-in-participants"),
			wantCode: http.StatusForbidden,
		},
		"with a valid cookie and a valid transaction": {
			transaction: app.SharedTransaction{
				Participants: []string{"user-01", "user-02"},
				Shop:         "test-shop",
				Amount:       123,
				Date:         123456,
			},
			cookie:   *CreateCookie("user-01"),
			wantCode: http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)
			request := app.NewCreateSharedTxnRequest(tc.transaction)
			request.AddCookie(&tc.cookie)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}

func addTransaction(h http.Handler, txn app.SharedTransaction) {
	request := app.NewCreateSharedTxnRequest(txn)
	request.AddCookie(CreateCookie(txn.Participants[0]))
	response := httptest.NewRecorder()
	h.ServeHTTP(response, request)
}

func TestGetTxnsByTracker(t *testing.T) {
	testTxn := app.SharedTransaction{
		Participants: []string{"user-01", "user-02"},
		Shop:         "test-shop",
		Amount:       123,
		Date:         123456,
		Tracker:      "test-tracker-01",
	}

	assert := assert.New(t)
	tests := map[string]struct {
		tracker          string
		wantTransactions []app.SharedTransaction
		wantCode         int
	}{
		"for a non-existent tracker or one with no transactions": {
			tracker:          "no-txn-tracker",
			wantTransactions: []app.SharedTransaction{},
			wantCode:         http.StatusOK,
		},
		"for a tracker with at least one transaction": {
			tracker:          "test-tracker-01",
			wantTransactions: []app.SharedTransaction{testTxn},
			wantCode:         http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)

			// add a transaction to be gotten
			addTransaction(router, testTxn)

			request := app.NewGetTxnsByTrackerRequest(tc.tracker)
			request.AddCookie(CreateCookie(TestSeanUser.ID))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)

			var got []app.SharedTransaction
			err := json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Fatalf("error parsing response from server %q into slice of shared txns: %v", response.Body, err)
			}
			assert.Len(got, len(tc.wantTransactions))

			// remove the ID from the got transactions to account for randomly generated
			var gotWithoutID []app.SharedTransaction
			for _, txn := range got {
				gotWithoutID = append(gotWithoutID, app.SharedTransaction{
					Participants: txn.Participants,
					Shop:         txn.Shop,
					Amount:       txn.Amount,
					Date:         txn.Date,
					Tracker:      txn.Tracker,
				})
			}

			assert.ElementsMatch(gotWithoutID, tc.wantTransactions)
		})
	}
}

var testUnsettledTxn = app.SharedTransaction{
	Participants: []string{"user-01", "user-02"},
	Shop:         "test-shop",
	Amount:       123,
	Date:         123456,
	Tracker:      "test-tracker-01",
	Unsettled:    true,
}

func TestGetUnsettledTxnsFromTracker(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		tracker          string
		wantTransactions []app.SharedTransaction
		wantCode         int
	}{
		"for a non-existent tracker or one with no unsettled txns": {
			tracker:          "no-unsettled-txn-tracker",
			wantTransactions: []app.SharedTransaction{},
			wantCode:         http.StatusOK,
		},
		"for a tracker with at least one unsettled txn": {
			tracker:          "test-tracker-01",
			wantTransactions: []app.SharedTransaction{testUnsettledTxn},
			wantCode:         http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)
			addTransaction(router, testUnsettledTxn)

			request := app.NewGetUnsettledTxnsByTrackerRequest(tc.tracker)
			request.AddCookie(CreateCookie(TestSeanUser.ID))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)

			var got []app.SharedTransaction
			err := json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Fatalf("error parsing response from server %q into slice of shared txns: %v", response.Body, err)
			}
			assert.Len(got, len(tc.wantTransactions))

			// remove the ID from the got transactions to account for randomly generated
			var gotWithoutID []app.SharedTransaction
			for _, txn := range got {
				gotWithoutID = append(gotWithoutID, app.SharedTransaction{
					Participants: txn.Participants,
					Shop:         txn.Shop,
					Amount:       txn.Amount,
					Date:         txn.Date,
					Tracker:      txn.Tracker,
					Unsettled:    txn.Unsettled,
				})
			}

			assert.ElementsMatch(gotWithoutID, tc.wantTransactions)
		})
	}
}

func TestSettleTxns(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]struct {
		cookie      http.Cookie
		initialTxns []app.SharedTransaction
		wantTxns    []app.SharedTransaction
		wantCode    int
	}{
		// "you cannot settle a transaction that you are not a participant of": {
		// 	cookie:      http.Cookie{Name: "session", Value: "not-a-participant"},
		// 	initialTxns: []app.SharedTransaction{testUnsettledTxn},
		// 	wantTxns:    []app.SharedTransaction{testUnsettledTxn},
		// 	wantCode:    http.StatusForbidden,
		// },
		"you can settle a transaction that you are a participant of": {
			cookie:      *CreateCookie(testUnsettledTxn.Participants[0]),
			initialTxns: []app.SharedTransaction{testUnsettledTxn},
			wantTxns:    []app.SharedTransaction{},
			wantCode:    http.StatusAccepted,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)
			for _, txn := range tc.initialTxns {
				addTransaction(router, txn)
			}

			// get the one unsettled transaction that should be there
			request := app.NewGetUnsettledTxnsByTrackerRequest(testUnsettledTxn.Tracker)
			request.AddCookie(&tc.cookie)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			var unsettled []app.SharedTransaction
			err := json.NewDecoder(response.Body).Decode(&unsettled)
			if err != nil {
				t.Fatalf("error parsing response from server %q into slice of shared txns: %v", response.Body, err)
			}

			request = app.NewSettleTxnsRequest(t, unsettled)
			request.AddCookie(&tc.cookie)
			response = httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)

			// get all of the requests from the unsettled txn tracker
			request = app.NewGetUnsettledTxnsByTrackerRequest(testUnsettledTxn.Tracker)
			request.AddCookie(&tc.cookie)
			response = httptest.NewRecorder()
			router.ServeHTTP(response, request)

			var got []app.SharedTransaction
			err = json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Fatalf("error parsing response from server %q into slice of shared txns: %v", response.Body, err)
			}
			assert.Len(got, len(tc.wantTxns))
		})
	}
}
