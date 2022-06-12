package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifahn/expenseus/internal/app"
	"github.com/stretchr/testify/assert"
)

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