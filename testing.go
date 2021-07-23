package expenseus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// NewGetExpenseByIDRequest creates a request to get an expense by id, adding the id to the request context.
func NewGetExpenseByIDRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/%s", id), nil)
	ctx := context.WithValue(req.Context(), "id", id)
	return req.WithContext(ctx)
}

func NewCreateExpenseRequest(user, name string) *http.Request {
	values := map[string]string{"user": user, "name": name}
	jsonValue, _ := json.Marshal(values)
	req, _ := http.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(jsonValue))
	return req
}

func NewGetExpenseByUserRequest(user string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/user/%s", user), nil)
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
