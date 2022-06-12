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

func TestUpdateTransactions(t *testing.T) {
	initialDetails := app.Transaction{
		Name:   "test-transaction",
		ID:     "test-id",
		UserID: TestSeanUser.ID,
		Amount: 100,
		Date:   333333,
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
				Name:   "new-name",
				UserID: initialDetails.UserID,
				Amount: 999,
				Date:   129384,
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
			createTestTransaction(t, router, test.initialTxnDetails, test.user.ID)

			// get all transactions to get the transaction that was just added
			request := app.NewGetAllTransactionsRequest()
			request.AddCookie(CreateCookie(test.user.ID))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			var transactionsGot []app.Transaction
			err := json.NewDecoder(response.Body).Decode(&transactionsGot)
			assert.NoError(err)
			assert.Len(transactionsGot, 1)

			if test.updateDetails.Name != "" {
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
