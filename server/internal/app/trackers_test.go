package app

import (
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
	request := NewCreateTrackerRequest(t, testTrackerDetails)
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(app.CreateTracker)
	handler.ServeHTTP(response, request)

	assert.Equal(http.StatusAccepted, response.Code)
	assert.Len(store.trackers, 1)
}

func TestGetTrackerByID(t *testing.T) {
	assert := assert.New(t)
	store := StubTransactionStore{
		trackers: []Tracker{TestTracker},
	}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	request := NewGetTrackerByIDRequest(t, TestTracker.ID)
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(app.GetTrackerByID)
	handler.ServeHTTP(response, request)

	var got Tracker
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from server %q into Tracker: '%v'", response.Body, err)
	}

	assert.Equal(jsonContentType, response.Result().Header.Get("content-type"))
	assert.Equal(http.StatusOK, response.Code)
	assert.Equal(TestTracker, got)
}

func TestGetTrackersByUser(t *testing.T) {
	assert := assert.New(t)
	store := StubTransactionStore{
		trackers: []Tracker{TestTracker},
	}
	app := New(&store, &StubOauthConfig{}, &StubSessionManager{}, "", &StubImageStore{})

	request := NewGetTrackerByUserRequest(t, TestSeanUser.ID)
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(app.GetTrackersByUser)
	handler.ServeHTTP(response, request)

	var got []Tracker
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("error parsing response from %q into slice of Trackers: '%v'", response.Body, err)
	}

	assert.Equal(jsonContentType, response.Result().Header.Get("content-type"))
	assert.Equal(http.StatusOK, response.Code)
	assert.Contains(got, TestTracker)
}
