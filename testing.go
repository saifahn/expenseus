package expenseus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

var TestSeanUser = User{
	Username: "saifahn",
	Name:     "Sean Li",
	ID:       "sean_id",
}

var TestTomomiUser = User{
	Username: "tomochi",
	Name:     "Tomomi Kinoshita",
	ID:       "tomomi_id",
}

var TestSeanExpense = Expense{
	ID:     "1",
	Name:   "Expense 1",
	UserID: TestSeanUser.ID,
}

var TestTomomiExpense = Expense{
	ID:     "9281",
	Name:   "Expense 9281",
	UserID: TestTomomiUser.ID,
}

var TestTomomiExpense2 = Expense{
	ID:     "14928",
	Name:   "Expense 14928",
	UserID: TestTomomiUser.ID,
}

// NewGetExpenseRequest creates a request to be used in tests get an expense
// by id, adding the id to the request context.
func NewGetExpenseRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/%s", id), nil)
	ctx := context.WithValue(req.Context(), CtxKeyExpenseID, id)
	return req.WithContext(ctx)
}

// NewCreateExpenseRequest creates a request to be used in tests to create an
// expense that is associated with a user.
func NewCreateExpenseRequest(user, name string) *http.Request {
	values := map[string]string{"user": user, "name": name}
	jsonValue, _ := json.Marshal(values)
	req, _ := http.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(jsonValue))
	return req
}

// NewGetExpensesByUsernameRequest creates a request to be used in tests to get all
// expenses of a user, adding the user to the request context.
func NewGetExpensesByUsernameRequest(username string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/user/%s", username), nil)
	ctx := context.WithValue(req.Context(), CtxKeyUsername, username)
	return req.WithContext(ctx)
}

func NewGetAllExpensesRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/expenses/", nil)
	return req
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got response body of %q, want %q", got, want)
	}
}

func AssertResponseStatus(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}
