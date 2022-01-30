package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTransactionByID(t *testing.T) {
	store := StubTransactionStore{
		users: []User{},
		transactions: map[string]Transaction{
			"1":    TestSeanTransaction,
			"9281": TestTomomiTransaction,
			"134":  TestTransactionWithImage,
		},
	}
	images := StubImageStore{}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &images)

	t.Run("get an transaction by id", func(t *testing.T) {
		request := NewGetTransactionRequest("1")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransaction)
		handler.ServeHTTP(response, request)

		var got Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Transaction, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, got, TestSeanTransaction)
	})

	t.Run("gets another transaction by id", func(t *testing.T) {
		request := NewGetTransactionRequest("9281")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransaction)
		handler.ServeHTTP(response, request)

		var got Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Transaction, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, got, TestTomomiTransaction)
	})

	t.Run("returns a response without an imageKey or imageUrl for an transaction without an image", func(t *testing.T) {
		request := NewGetTransactionRequest("9281")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransaction)
		handler.ServeHTTP(response, request)

		rawJSON := response.Body.String()

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.NotContains(t, rawJSON, "imageKey")
		assert.NotContains(t, rawJSON, "imageUrl")
	})

	t.Run("returns a response with an imageUrl for an transaction that has an image", func(t *testing.T) {
		request := NewGetTransactionRequest("134")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransaction)
		handler.ServeHTTP(response, request)

		rawJSON := response.Body.String()

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		// TODO somehow don't return this to the front end but we need to use it in the back end so we can't just use "-"?
		// assert.NotContains(t, rawJSON, "imageKey")
		assert.Len(t, images.addImageToTransactionCalls, 1)
		assert.Contains(t, rawJSON, "imageUrl")
	})

	t.Run("returns 404 on non-existent transaction", func(t *testing.T) {
		request := NewGetTransactionRequest("13371337")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransaction)
		handler.ServeHTTP(response, request)

		var got Transaction
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Transaction, '%v'", response.Body, err)
		}

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

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
		request := NewGetTransactionsByUsernameRequest(TestTomomiUser.Username)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransactionsByUsername)
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
		request := NewGetTransactionsByUsernameRequest(TestSeanUser.Username)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetTransactionsByUsername)
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
	t.Run("returns an error if there is no user in the context", func(t *testing.T) {
		store := StubTransactionStore{}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		values := map[string]io.Reader{
			"expenseName": strings.NewReader("Test Transaction"),
		}

		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		for _, reader := range values {
			var fw io.Writer
			fw, _ = w.CreateFormField("expenseName")
			if _, err := io.Copy(fw, reader); err != nil {
				fmt.Println(err.Error())
			}
		}
		w.Close()
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/expenses", &b)
		request.Header.Set("Content-Type", w.FormDataContentType())
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateTransaction)
		assert.Panics(t, func() {
			handler.ServeHTTP(response, request)
		}, "The code did not panic due to a lack of context")
	})

	t.Run("creates a new transaction on POST", func(t *testing.T) {
		store := StubTransactionStore{
			transactions: map[string]Transaction{},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		values := map[string]io.Reader{
			"expenseName": strings.NewReader("Test Transaction"),
		}
		request := addUserCookieAndContext(NewCreateTransactionRequest(values), TestTomomiUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateTransaction)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
		assert.Len(t, store.recordTransactionCalls, 1)
	})

	// prepares a temp file, information, and values for image upload tests
	prepareFileAndInfo := func(t *testing.T) (*os.File, string, map[string]io.Reader) {
		f, err := os.CreateTemp("", "example-file")
		if err != nil {
			t.Fatal(err)
		}
		transactionName := "Test Transaction with Image"

		values := map[string]io.Reader{
			"expenseName": strings.NewReader(transactionName),
			"image":       f,
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

		request := addUserCookieAndContext(NewCreateTransactionRequest(values), TestSeanUser.ID)
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

		request := addUserCookieAndContext(NewCreateTransactionRequest(values), userID)
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
