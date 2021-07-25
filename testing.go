package expenseus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// NewGetExpenseByIDRequest creates a request to be used in tests get an expense
// by id, adding the id to the request context.
func NewGetExpenseByIDRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/%s", id), nil)
	ctx := context.WithValue(req.Context(), "expenseID", id)
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

// NewGetExpensesByUserRequest creates a request to be used in tests to get all
// expenses of a user, adding the user to the request context.
func NewGetExpensesByUserRequest(username string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/user/%s", username), nil)
	ctx := context.WithValue(req.Context(), "username", username)
	return req.WithContext(ctx)
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
