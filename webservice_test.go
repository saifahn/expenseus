package expenseus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExpenses(t *testing.T) {
	store := StubExpenseStore{
		map[string]string{
			"1":    "Expense 1",
			"9281": "Expense 9281",
		},
	}
	webservice := &WebService{&store}

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

	t.Run("returns 404 on non-existent expense", func(t *testing.T) {
		request := newGetExpenseRequest("13371337")
		response := httptest.NewRecorder()

		webservice.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		if got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
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

type StubExpenseStore struct {
	expenses map[string]string
}

func (s *StubExpenseStore) GetExpense(id string) (expense string) {
	return s.expenses[id]
}
