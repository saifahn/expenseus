package app_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testUser := app.User{
		Username: "test-user",
	}

	expectFn := func(ma *mock_app.App) {
		ma.MockStore.EXPECT().CreateUser(testUser).Return(nil).Times(1)
	}

	assert := assert.New(t)
	a := mock_app.SetUp(t, expectFn)

	userJSON, err := json.Marshal(testUser)
	assert.NoError(err)

	request := app.NewCreateUserRequest(userJSON)
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(a.CreateUser)
	handler.ServeHTTP(response, request)

	assert.Equal(http.StatusAccepted, response.Code)
}

func TestListUsers(t *testing.T) {
	testUser := app.User{
		Username: "test-user",
	}

	tests := map[string]struct {
		users    []app.User
		expectFn mock_app.MockAppFn
		wantCode int
	}{
		"with no users returned from the store": {
			users: []app.User{},
			expectFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetAllUsers().Return([]app.User{}, nil).Times(1)
			},
			wantCode: http.StatusOK,
		},
		"with a user returned from the store": {
			users: []app.User{testUser},
			expectFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetAllUsers().Return([]app.User{testUser}, nil).Times(1)
			},
			wantCode: http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectFn)

			request := app.NewGetAllUsersRequest()
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.ListUsers)
			handler.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)

			var got []app.User
			err := json.NewDecoder(response.Body).Decode(&got)
			assert.NoError(err)
			assert.Equal(tc.users, got)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	testUser := app.User{
		Username: "test-user",
	}

	tests := map[string]struct {
		userID   string
		user     app.User
		expectFn mock_app.MockAppFn
		wantCode int
	}{
		"with no user returned from the store": {
			userID: "not-found",
			user:   app.User{},
			expectFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetUser("not-found").Return(app.User{}, app.ErrDBItemNotFound).Times(1)
			},
			wantCode: http.StatusNotFound,
		},
		"with a user returned from the store": {
			userID: testUser.Username,
			user:   testUser,
			expectFn: func(ma *mock_app.App) {
				ma.MockStore.EXPECT().GetUser(testUser.Username).Return(testUser, nil).Times(1)
			},
			wantCode: http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectFn)

			request := app.NewGetUserRequest(tc.userID)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.GetUser)
			handler.ServeHTTP(response, request)

			assert.Equal(tc.wantCode, response.Code)

			if tc.wantCode == http.StatusOK {
				var got app.User
				err := json.NewDecoder(response.Body).Decode(&got)
				assert.NoError(err)
				assert.Equal(tc.user, got)
			}
		})
	}
}
