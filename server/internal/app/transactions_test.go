package app

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTransactionByUser(t *testing.T) {
	store := StubTransactionStore{
		users: []User{
			TestSeanUser,
			TestTomomiUser,
		},
		transactions: map[string]Transaction{
			"1":    TestSeanTransaction,
			"9281": TestTomomiTransaction,
		},
	}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	t.Run("gets tomochi's transactions", func(t *testing.T) {
		request := NewGetTransactionsByUserRequest(TestTomomiUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransactionsByUser)
		handler.ServeHTTP(response, request)

		var got []Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Len(t, got, 1)
		assert.Contains(t, got, TestTomomiTransaction)
	})

	t.Run("gets saifahn's transactions", func(t *testing.T) {
		request := NewGetTransactionsByUserRequest(TestSeanUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransactionsByUser)
		handler.ServeHTTP(response, request)

		var got []Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Len(t, got, 1)
		assert.Contains(t, got, TestSeanTransaction)
	})
}

func TestCreateTransaction(t *testing.T) {
	// FIXME: for some reason, this stopped working after adding amount to the transactions
	// t.Run("returns an error if there is no user in the context", func(t *testing.T) {
	// 	store := StubTransactionStore{}
	// 	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	// 	values := map[string]io.Reader{
	// 		"transactionName": strings.NewReader("Test Transaction"),
	// 		"amount":          strings.NewReader("123"),
	// 	}

	// 	var b bytes.Buffer
	// 	w := multipart.NewWriter(&b)
	// 	for _, reader := range values {
	// 		var fw io.Writer
	// 		fw, _ = w.CreateFormField("transactionName")
	// 		if _, err := io.Copy(fw, reader); err != nil {
	// 			fmt.Println(err.Error())
	// 		}
	// 	}
	// 	w.Close()
	// 	request, _ := http.NewRequest(http.MethodPost, "/api/v1/transactions", &b)
	// 	request.Header.Set("Content-Type", w.FormDataContentType())
	// 	response := httptest.NewRecorder()

	// 	handler := http.HandlerFunc(app.CreateTransaction)
	// 	assert.Panics(t, func() {
	// 		handler.ServeHTTP(response, request)
	// 	}, "The code did not panic due to a lack of context")
	// })

	t.Run("calls the store's CreateTransaction with the right details on POST", func(t *testing.T) {
		store := StubTransactionStore{
			transactions: map[string]Transaction{},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})
		expectedDetails := TransactionDetails{
			Name:   "Test Transaction",
			Amount: 123,
			Date:   1644085875,
			UserID: TestTomomiUser.ID,
		}
		payload := MakeCreateTransactionRequestPayload(expectedDetails)
		request := AddUserCookieAndContext(NewCreateTransactionRequest(payload), TestTomomiUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateTransaction)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
		assert.Len(t, store.recordTransactionCalls, 1)
		assert.Equal(t, expectedDetails, store.recordTransactionCalls[0])
	})

	// prepares a temp file, information, and values for image upload tests
	var testAmount = "918"
	var testAmountInt64 = int64(918)
	var testDate = "1644085875"
	var testDateInt64 = int64(1644085875)
	prepareFileAndInfo := func(t *testing.T) (*os.File, string, map[string]io.Reader) {
		f, err := os.CreateTemp("", "example-file")
		if err != nil {
			t.Fatal(err)
		}
		transactionName := "Test Transaction with Image"

		values := map[string]io.Reader{
			"transactionName": strings.NewReader(transactionName),
			"amount":          strings.NewReader(testAmount),
			"date":            strings.NewReader(testDate),
			"image":           f,
		}
		return f, transactionName, values
	}

	t.Run("if an image is provided and it fails the image check, there is an error response", func(t *testing.T) {
		store := StubTransactionStore{
			transactions: map[string]Transaction{},
		}
		images := StubInvalidImageStore{}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &images)

		f, _, values := prepareFileAndInfo(t)
		defer f.Close()
		defer os.Remove(f.Name())

		request := AddUserCookieAndContext(NewCreateTransactionRequest(values), TestSeanUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateTransaction)
		handler.ServeHTTP(response, request)

		// the invalid image store will return this error if the image is invalid
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Len(t, images.uploadCalls, 0)
	})

	t.Run("if an image is provided and the image check is successful, the image is uploaded and an transaction is created with an image key", func(t *testing.T) {
		store := StubTransactionStore{
			transactions: map[string]Transaction{},
		}
		images := StubImageStore{}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &images)
		userID := TestSeanUser.ID

		f, transactionName, values := prepareFileAndInfo(t)
		defer f.Close()
		defer os.Remove(f.Name())

		request := AddUserCookieAndContext(NewCreateTransactionRequest(values), userID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateTransaction)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
		assert.Len(t, images.uploadCalls, 1)
		assert.Len(t, store.recordTransactionCalls, 1)
		got := store.recordTransactionCalls[0]
		want := TransactionDetails{
			Name:     transactionName,
			UserID:   userID,
			ImageKey: testImageKey,
			Amount:   testAmountInt64,
			Date:     testDateInt64,
		}
		assert.Equal(t, want, got)
	})
}

func TestGetAllTransactions(t *testing.T) {
	t.Run("gets all transactions with one transaction", func(t *testing.T) {
		wantedTransactions := []Transaction{
			TestTomomiTransaction,
		}
		store := StubTransactionStore{
			users: []User{},
			transactions: map[string]Transaction{
				"9281": TestTomomiTransaction,
			},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetAllTransactionsRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetAllTransactions)
		handler.ServeHTTP(response, request)

		var got []Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(wantedTransactions), len(got))
		assert.ElementsMatch(t, got, wantedTransactions)
	})

	t.Run("gets all transactions with more than one transaction", func(t *testing.T) {
		wantedTransactions := []Transaction{
			TestSeanTransaction, TestTomomiTransaction, TestTomomiTransaction2,
		}
		store := StubTransactionStore{
			users: []User{},
			transactions: map[string]Transaction{
				"1":     TestSeanTransaction,
				"9281":  TestTomomiTransaction,
				"14928": TestTomomiTransaction2,
			},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetAllTransactionsRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetAllTransactions)
		handler.ServeHTTP(response, request)

		var got []Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Transactions, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(wantedTransactions), len(got))
		assert.ElementsMatch(t, got, wantedTransactions)
	})
}
