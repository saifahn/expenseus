package expenseus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExpenses(t *testing.T) {
	webservice := &WebService{}
	t.Run("get an expense by id", func(t *testing.T) {
		request := newGetExpenseRequest("1")
		response := httptest.NewRecorder()

		webservice.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "Expense 1")
	})

	t.Run("gets another expense by id", func(t *testing.T) {
		request := newGetExpenseRequest("9281")
		response := httptest.NewRecorder()

		webservice.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "Expense 9281")
	})
}

func newGetExpenseRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/%s", id), nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

}
