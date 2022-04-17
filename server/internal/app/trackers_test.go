package app

import (
	"bytes"
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
