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

func TestGetExpenseByID(t *testing.T) {
	store := StubExpenseStore{
		users: []User{},
		expenses: map[string]Expense{
			"1":    TestSeanExpense,
			"9281": TestTomomiExpense,
			"134":  TestExpenseWithImage,
		},
	}
	images := StubImageStore{}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &images)

	t.Run("get an expense by id", func(t *testing.T) {
		request := NewGetExpenseRequest("1")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, got, TestSeanExpense)
	})

	t.Run("gets another expense by id", func(t *testing.T) {
		request := NewGetExpenseRequest("9281")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, got, TestTomomiExpense)
	})

	t.Run("returns a response without an imageKey or imageUrl for an expense without an image", func(t *testing.T) {
		request := NewGetExpenseRequest("9281")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetExpense)
		handler.ServeHTTP(response, request)

		rawJSON := response.Body.String()

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.NotContains(t, rawJSON, "imageKey")
		assert.NotContains(t, rawJSON, "imageUrl")
	})

	t.Run("returns a response with an imageUrl for an expense that has an image", func(t *testing.T) {
		request := NewGetExpenseRequest("134")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetExpense)
		handler.ServeHTTP(response, request)

		rawJSON := response.Body.String()

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		// TODO somehow don't return this to the front end but we need to use it in the back end so we can't just use "-"?
		// assert.NotContains(t, rawJSON, "imageKey")
		assert.Len(t, images.addImageToExpenseCalls, 1)
		assert.Contains(t, rawJSON, "imageUrl")
	})

	t.Run("returns 404 on non-existent expense", func(t *testing.T) {
		request := NewGetExpenseRequest("13371337")
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetExpense)
		handler.ServeHTTP(response, request)

		var got Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into Expense, '%v'", response.Body, err)
		}

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetExpenseByUser(t *testing.T) {
	store := StubExpenseStore{
		users: []User{
			TestSeanUser,
			TestTomomiUser,
		},
		expenses: map[string]Expense{
			"1":    TestSeanExpense,
			"9281": TestTomomiExpense,
		},
	}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	t.Run("gets tomochi's expenses", func(t *testing.T) {
		request := NewGetExpensesByUsernameRequest(TestTomomiUser.Username)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetExpensesByUsername)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Len(t, got, 1)
		assert.Contains(t, got, TestTomomiExpense)
	})

	t.Run("gets saifahn's expenses", func(t *testing.T) {
		request := NewGetExpensesByUsernameRequest(TestSeanUser.Username)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetExpensesByUsername)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Len(t, got, 1)
		assert.Contains(t, got, TestSeanExpense)
	})
}

func TestCreateExpense(t *testing.T) {
	t.Run("returns an error if there is no user in the context", func(t *testing.T) {
		store := StubExpenseStore{}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		values := map[string]io.Reader{
			"expenseName": strings.NewReader("Test Expense"),
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

		handler := http.HandlerFunc(app.CreateExpense)
		assert.Panics(t, func() {
			handler.ServeHTTP(response, request)
		}, "The code did not panic due to a lack of context")
	})

	t.Run("creates a new expense on POST", func(t *testing.T) {
		store := StubExpenseStore{
			expenses: map[string]Expense{},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		values := map[string]io.Reader{
			"expenseName": strings.NewReader("Test Expense"),
		}
		request := addUserCookieAndContext(NewCreateExpenseRequest(values), TestTomomiUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateExpense)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
		assert.Len(t, store.recordExpenseCalls, 1)
	})

	// prepares a temp file, information, and values for image upload tests
	prepareFileAndInfo := func(t *testing.T) (*os.File, string, map[string]io.Reader) {
		f, err := os.CreateTemp("", "example-file")
		if err != nil {
			t.Fatal(err)
		}
		expenseName := "Test Expense with Image"

		values := map[string]io.Reader{
			"expenseName": strings.NewReader(expenseName),
			"image":       f,
		}
		return f, expenseName, values
	}

	t.Run("if an image is provided and it fails the image check, there is an error response", func(t *testing.T) {
		store := StubExpenseStore{
			expenses: map[string]Expense{},
		}
		images := StubInvalidImageStore{}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &images)

		f, _, values := prepareFileAndInfo(t)
		defer f.Close()
		defer os.Remove(f.Name())

		request := addUserCookieAndContext(NewCreateExpenseRequest(values), TestSeanUser.ID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateExpense)
		handler.ServeHTTP(response, request)

		// the invalid image store will return this error if the image is invalid
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
		assert.Len(t, images.uploadCalls, 0)
	})

	t.Run("if an image is provided and the image check is successful, the image is uploaded and an expense is created with an image key", func(t *testing.T) {
		store := StubExpenseStore{
			expenses: map[string]Expense{},
		}
		images := StubImageStore{}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &images)
		userID := TestSeanUser.ID

		f, expenseName, values := prepareFileAndInfo(t)
		defer f.Close()
		defer os.Remove(f.Name())

		request := addUserCookieAndContext(NewCreateExpenseRequest(values), userID)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.CreateExpense)
		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
		assert.Len(t, images.uploadCalls, 1)
		assert.Len(t, store.recordExpenseCalls, 1)
		got := store.recordExpenseCalls[0]
		want := ExpenseDetails{
			Name:     expenseName,
			UserID:   userID,
			ImageKey: testImageKey,
		}
		assert.Equal(t, want, got)
	})
}

func TestGetAllExpenses(t *testing.T) {
	t.Run("gets all expenses with one expense", func(t *testing.T) {
		wantedExpenses := []Expense{
			TestTomomiExpense,
		}
		store := StubExpenseStore{
			users: []User{},
			expenses: map[string]Expense{
				"9281": TestTomomiExpense,
			},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetAllExpensesRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetAllExpenses)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(wantedExpenses), len(got))
		assert.ElementsMatch(t, got, wantedExpenses)
	})

	t.Run("gets all expenses with more than one expense", func(t *testing.T) {
		wantedExpenses := []Expense{
			TestSeanExpense, TestTomomiExpense, TestTomomiExpense2,
		}
		store := StubExpenseStore{
			users: []User{},
			expenses: map[string]Expense{
				"1":     TestSeanExpense,
				"9281":  TestTomomiExpense,
				"14928": TestTomomiExpense2,
			},
		}
		app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

		request := NewGetAllExpensesRequest()
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.GetAllExpenses)
		handler.ServeHTTP(response, request)

		var got []Expense
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("error parsing response from server %q into slice of Expenses, '%v'", response.Body, err)
		}

		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(wantedExpenses), len(got))
		assert.ElementsMatch(t, got, wantedExpenses)
	})

}
