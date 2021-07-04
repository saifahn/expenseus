package expenseus

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExpenses(t *testing.T) {
	t.Run("get an expense by id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/expenses/1", nil)
		response := httptest.NewRecorder()

		WebService(response, request)

		got := response.Body.String()
		want := "Expense 1"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("gets another expense by id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/expenses/9281", nil)
		response := httptest.NewRecorder()

		WebService(response, request)

		got := response.Body.String()
		want := "Expense 9281"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
