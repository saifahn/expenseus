package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifahn/expenseus/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestCreatingTransactionsAndRetrievingThem(t *testing.T) {
	t.Run("an transaction can be added with a valid cookie and be retrieved as part of a GetAll request", func(t *testing.T) {
		router, tearDownDB := SetUpTestServer(t)
		defer tearDownDB(t)
		assert := assert.New(t)

		// create user in the db
		CreateUser(t, TestSeanUser, router)

		// create a transaction and store it
		wantTxnDetails := TestSeanTxnDetails
		CreateTestTxn(t, router, wantTxnDetails, wantTxnDetails.UserID)

		// try and get it
		request := app.NewGetTransactionsByUserRequest(wantTxnDetails.UserID)
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
		CreateTestTxn(t, router, wantTxnDetails, wantTxnDetails.UserID)

		request := app.NewGetTransactionsByUserRequest(wantTxnDetails.UserID)
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
}

func TestCreatingTxns(t *testing.T) {
	txnWithoutCategory := app.Transaction{
		Location: "test-location",
		UserID:   TestSeanUser.ID,
		Amount:   200,
		Date:     8972813,
		Details:  "without-category",
	}

	txnWithCategory := app.Transaction{
		Location: "test-location",
		UserID:   TestSeanUser.ID,
		Amount:   200,
		Date:     8972813,
		Category: "other.other",
		Details:  "with-category",
	}

	tests := map[string]struct {
		txnDetails app.Transaction
		wantTxns   []app.Transaction
		wantCode   int
	}{
		"with a txn that has no category": {
			txnDetails: txnWithoutCategory,
			wantTxns:   []app.Transaction{},
			wantCode:   http.StatusBadRequest,
		},
		"with a txn that has a category": {
			txnDetails: txnWithCategory,
			wantTxns:   []app.Transaction{txnWithCategory},
			wantCode:   http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)
			assert := assert.New(t)

			CreateUser(t, TestSeanUser, router)
			CreateTestTxn(t, router, tc.txnDetails, TestSeanUser.ID)
			request := app.NewGetTransactionsByUserRequest(TestSeanUser.ID)
			request.AddCookie(CreateCookie(TestSeanUser.ID))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(http.StatusOK, response.Code)

			if len(tc.wantTxns) != 0 {
				var transactionsGot []app.Transaction
				err := json.NewDecoder(response.Body).Decode(&transactionsGot)
				assert.NoError(err)
				assert.Len(transactionsGot, len(tc.wantTxns))
				AssertEqualTxnDetails(t, tc.wantTxns[0], transactionsGot[0])
			}
		})
	}
}

func TestGetTxnsByUser(t *testing.T) {
	initTxn1 := app.Transaction{
		Location: "test-location",
		UserID:   "a-user",
		Amount:   300,
		Date:     10000,
		Category: "something",
	}

	tests := map[string]struct {
		initTxns []app.Transaction
		wantTxns []app.Transaction
		user     string
		wantCode int
	}{
		"with a user that has no transactions": {
			initTxns: []app.Transaction{initTxn1},
			wantTxns: []app.Transaction{},
			user:     "no-txns-user",
			wantCode: http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)
			assert := assert.New(t)

			CreateUser(t, TestSeanUser, router)
			for _, txn := range tc.initTxns {
				CreateTestTxn(t, router, txn, txn.ID)
			}

			request := app.NewGetTransactionsByUserRequest(tc.user)
			request.AddCookie(CreateCookie(tc.user))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)

			var got []app.Transaction
			err := json.NewDecoder(response.Body).Decode(&got)
			assert.NoError(err)
			assert.Len(got, len(tc.wantTxns))
		})
	}
}

func TestGetTxnBetweenDates(t *testing.T) {
	initialTxn := app.Transaction{
		Location: "test-location",
		UserID:   "a-user",
		Amount:   300,
		Date:     10000,
		Category: "something",
	}

	tests := map[string]struct {
		wantTxns []app.Transaction
		from     int64
		to       int64
		wantCode int
	}{
		"with a date range containing a transaction": {
			wantTxns: []app.Transaction{initialTxn},
			from:     10000,
			to:       11000,
			wantCode: http.StatusOK,
		},
		"with a date range not containing a transaction": {
			wantTxns: []app.Transaction{},
			from:     20000,
			to:       21000,
			wantCode: http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)
			assert := assert.New(t)

			CreateUser(t, TestSeanUser, router)
			CreateTestTxn(t, router, initialTxn, initialTxn.UserID)

			request := app.NewGetTxnsBetweenDatesRequest(initialTxn.UserID, tc.from, tc.to)
			request.AddCookie(CreateCookie(initialTxn.UserID))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)

			var got []app.Transaction
			err := json.NewDecoder(response.Body).Decode(&got)
			assert.NoError(err)
			assert.Len(got, len(tc.wantTxns))

			if len(tc.wantTxns) > 0 {
				for i := range tc.wantTxns {
					AssertEqualTxnDetails(t, got[i], tc.wantTxns[i])
				}
			}
		})
	}
}

func TestUpdateTransactions(t *testing.T) {
	initialDetails := app.Transaction{
		Location: "test-location",
		ID:       "test-id",
		UserID:   TestSeanUser.ID,
		Amount:   100,
		Date:     333333,
		Category: "test.category",
		Details:  "test-details",
	}
	tests := map[string]struct {
		initialTxnDetails app.Transaction
		updateDetails     app.Transaction
		user              app.User
		wantCode          int
	}{
		"attempting to update a non-existent transaction returns a 404": {
			initialTxnDetails: initialDetails,
			updateDetails:     app.Transaction{},
			user:              TestSeanUser,
			wantCode:          http.StatusNotFound,
		},
		"an existing transaction can be updated": {
			initialTxnDetails: initialDetails,
			updateDetails: app.Transaction{
				Location: "new-location",
				UserID:   initialDetails.UserID,
				Amount:   999,
				Date:     129384,
				Category: "new.category",
				Details:  "new-details",
			},
			user:     TestSeanUser,
			wantCode: http.StatusAccepted,
		},
		// TODO: don't allow a different user to update the details
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			router, tearDownDB := SetUpTestServer(t)
			defer tearDownDB(t)
			assert := assert.New(t)

			CreateUser(t, test.user, router)
			CreateTestTxn(t, router, test.initialTxnDetails, test.user.ID)

			// get transactions to get the transaction that was just added
			request := app.NewGetTransactionsByUserRequest(test.user.ID)
			request.AddCookie(CreateCookie(test.user.ID))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			var transactionsGot []app.Transaction
			err := json.NewDecoder(response.Body).Decode(&transactionsGot)
			assert.NoError(err)
			assert.Len(transactionsGot, 1)

			if test.updateDetails.Location != "" {
				test.updateDetails.ID = transactionsGot[0].ID
			}
			// update the transaction
			request = app.NewUpdateTransactionRequest(test.updateDetails)
			request.AddCookie(CreateCookie(test.user.ID))
			response = httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(test.wantCode, response.Code)

			// TODO: potentially get the initial transaction and see if it matches tc.wantTxnDetails
		})
	}
}

func TestUpdateThenGetByDateRange(t *testing.T) {
	initialDetails := app.Transaction{
		Location: "test-location",
		UserID:   TestSeanUser.ID,
		Amount:   100,
		Date:     333333,
		Category: "test.category",
		Details:  "test-details",
	}
	// create a txn
	router, tearDownDB := SetUpTestServer(t)
	defer tearDownDB(t)
	assert := assert.New(t)

	CreateUser(t, TestSeanUser, router)
	CreateTestTxn(t, router, initialDetails, TestSeanUser.ID)

	request := app.NewGetTxnsBetweenDatesRequest(TestSeanUser.ID, 30000, 40000)
	request.AddCookie(CreateCookie(TestSeanUser.ID))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var transactionsGot []app.Transaction
	err := json.NewDecoder(response.Body).Decode(&transactionsGot)
	assert.NoError(err)
	assert.Len(transactionsGot, 1)
	AssertEqualTxnDetails(t, initialDetails, transactionsGot[0])

	// update the transaction date
	updatedDetails := initialDetails
	updatedDetails.ID = transactionsGot[0].ID
	updatedDetails.Date = 50000
	request = app.NewUpdateTransactionRequest(updatedDetails)
	request.AddCookie(CreateCookie(TestSeanUser.ID))
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// get the new one with a new date range
	request = app.NewGetTxnsBetweenDatesRequest(TestSeanUser.ID, 40000, 50000)
	request.AddCookie(CreateCookie(TestSeanUser.ID))
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	transactionsGot = []app.Transaction{}
	err = json.NewDecoder(response.Body).Decode(&transactionsGot)
	assert.NoError(err)
	assert.Len(transactionsGot, 1)
	AssertEqualTxnDetails(t, updatedDetails, transactionsGot[0])
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
			CreateTestTxn(t, router, tc.td, tc.user)

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
