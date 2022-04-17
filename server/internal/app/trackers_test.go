package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTracker(t *testing.T) {
	assert := assert.New(t)
	store := StubTransactionStore{}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	testTrackerDetails := Tracker{
		Name:  "Test Tracker",
		Users: []string{TestSeanUser.ID},
	}
	trackerJSON, err := json.Marshal(testTrackerDetails)
	if err != nil {
		t.Fatalf(err.Error())
	}

	request, _ := http.NewRequest(http.MethodPost, "/api/v1/trackers", bytes.NewBuffer(trackerJSON))
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(app.CreateTracker)
	handler.ServeHTTP(response, request)

	assert.Equal(http.StatusAccepted, response.Code)
	assert.Len(store.trackers, 1)
}

func TestGetTrackerByID(t *testing.T) {
	assert := assert.New(t)
	testTracker := Tracker{
		Name:  "Test Tracker",
		Users: []string{TestSeanUser.ID},
		ID:    "test-id",
	}
	store := StubTransactionStore{
		trackers: []Tracker{testTracker},
	}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	request, _ := http.NewRequest(http.MethodGet, "/api/v1/trackers", nil)
	request = request.WithContext(context.WithValue(request.Context(), CtxKeyTrackerID, testTracker.ID))
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(app.GetTrackerByID)
	handler.ServeHTTP(response, request)

	var got Tracker
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into Tracker, '%v'", response.Body, err)
	}

	assert.Equal(jsonContentType, response.Result().Header.Get("content-type"))
	assert.Equal(http.StatusOK, response.Code)
	assert.Equal(testTracker, got)
}
