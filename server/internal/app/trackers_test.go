package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

var testTracker = app.Tracker{
	Name:  "test-tracker",
	Users: []string{"test-user"},
}

func TestCreateTracker(t *testing.T) {
	expectFn := func(ma *mock_app.App) {
		ma.MockStore.EXPECT().CreateTracker(testTracker).Return(nil).Times(1)
	}

	assert := assert.New(t)
	a := mock_app.SetUp(t, expectFn)

	req := app.NewCreateTrackerRequest(t, testTracker)
	req = req.WithContext(context.WithValue(req.Context(), app.CtxKeyUserID, "test-user"))
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(a.CreateTracker)
	handler.ServeHTTP(response, req)

	assert.Equal(http.StatusAccepted, response.Code)
}

func TestGetTrackerByID(t *testing.T) {
	expectFn := func(ma *mock_app.App) {
		ma.MockStore.EXPECT().GetTracker(testTracker.ID).Return(testTracker, nil).Times(1)
	}

	assert := assert.New(t)
	a := mock_app.SetUp(t, expectFn)

	req := app.NewGetTrackerByIDRequest(testTracker.ID)
	req = req.WithContext(context.WithValue(req.Context(), app.CtxKeyUserID, "test-user"))
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(a.GetTrackerByID)
	handler.ServeHTTP(response, req)

	assert.Equal(http.StatusOK, response.Code)
	// TODO: table test for non-existent
}

func TestGetTrackersByUser(t *testing.T) {
	expectFn := func(ma *mock_app.App) {
		ma.MockStore.EXPECT().GetTrackersByUser("test-user").Return([]app.Tracker{testTracker}, nil).Times(1)
	}

	assert := assert.New(t)
	a := mock_app.SetUp(t, expectFn)

	request := app.NewGetTrackerByUserRequest("test-user")
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(a.GetTrackersByUser)
	handler.ServeHTTP(response, request)

	assert.Equal(http.StatusOK, response.Code)
	// TODO: table test for non-existent
}
