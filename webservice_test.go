package expenseus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testSeanExpense = Expense{
	name: "Expense 1",
	user: "sean",
}

var testTomomiExpense = Expense{
	name: "Expense 9281",
	user: "tomomi",
}

func TestGetExpenseById(t *testing.T) {
	store := StubExpenseStore{
		map[string]Expense{
			"1":    testSeanExpense,
			"9281": testTomomiExpense,
		},
	}
	webservice := &WebService{&store}

	t.Run("get an expense by id", func(t *testing.T) {
		request := newGetExpenseByIdRequest("1")
		response := httptest.NewRecorder()

		webservice.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), testSeanExpense.name)
	})

	t.Run("gets another expense by id", func(t *testing.T) {
		request := newGetExpenseByIdRequest("9281")
		response := httptest.NewRecorder()

		webservice.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), testTomomiExpense.name)
	})

	t.Run("returns 404 on non-existent expense", func(t *testing.T) {
		request := newGetExpenseByIdRequest("13371337")
		response := httptest.NewRecorder()

		webservice.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestGetExpenseByUser(t *testing.T) {
	store := StubExpenseStore{
		map[string]Expense{
			"1":    testSeanExpense,
			"9281": testTomomiExpense,
		},
	}
	webservice := NewWebService(&store)

	t.Run("gets Tomomi's expenses", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/expenses/user/tomomi", nil)
		response := httptest.NewRecorder()

		webservice.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), fmt.Sprintf("[%v]", testTomomiExpense.name))
	})

	t.Run("gets Sean's expenses", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/expenses/user/sean", nil)
		response := httptest.NewRecorder()

		webservice.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), fmt.Sprintf("[%v]", testSeanExpense.name))
	})
}

func newGetExpenseByIdRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/%s", id), nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got response body of %q, want %q", got, want)
	}
}

func assertResponseStatus(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

type StubExpenseStore struct {
	expenses map[string]Expense
}

func (s *StubExpenseStore) GetExpenseNameById(id string) string {
	expense := s.expenses[id]
	return expense.name
}

func (s *StubExpenseStore) GetExpenseNamesByUser(user string) []string {
	var expenseNames []string
	for _, e := range s.expenses {
		if e.user == user {
			expenseNames = append(expenseNames, e.name)
		}
	}
	return expenseNames
}
