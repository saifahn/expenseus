package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/saifahn/expenseus/internal/app"
	mock_app "github.com/saifahn/expenseus/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestOauthLogin(t *testing.T) {
	assert := assert.New(t)
	oauthProviderURL := "test-oauth-url"

	expectFn := func(ma *mock_app.App) {
		ma.MockAuth.EXPECT().AuthCodeURL(gomock.Any()).Return(oauthProviderURL).Times(1)
	}
	a := mock_app.SetUp(t, expectFn)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/login_google", nil)
	response := httptest.NewRecorder()

	handler := http.HandlerFunc(a.OauthLogin)
	handler.ServeHTTP(response, req)

	assert.Equal(http.StatusTemporaryRedirect, response.Code)
	wantURL := fmt.Sprintf("/api/v1/%s", oauthProviderURL)
	assert.Equal(wantURL, response.Header().Get("Location"))
}

func TestOauthCallback(t *testing.T) {
	newUser := app.User{
		Username: "a-new-user",
		ID:       "a-new-user-id",
	}

	tests := map[string]struct {
		expectFn mock_app.MockAppFn
		wantCode int
	}{
		"when the user does not exist yet in the db": {
			expectFn: func(ma *mock_app.App) {
				ma.MockAuth.EXPECT().GetInfoAndGenerateUser(gomock.Any(), gomock.Any()).Return(newUser, nil).Times(1)
				ma.MockStore.EXPECT().GetUser(newUser.ID).Return(app.User{}, app.ErrDBItemNotFound).Times(1)
				ma.MockStore.EXPECT().CreateUser(newUser).Times(1)
				ma.MockSessions.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1)
			},
			wantCode: http.StatusTemporaryRedirect,
		},
		"when the user exists in the db": {
			expectFn: func(ma *mock_app.App) {
				ma.MockAuth.EXPECT().GetInfoAndGenerateUser(gomock.Any(), gomock.Any()).Return(newUser, nil).Times(1)
				ma.MockStore.EXPECT().GetUser(newUser.ID).Return(app.User{Username: newUser.Username, ID: newUser.ID}, nil).Times(1)
				ma.MockSessions.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1)
			},
			wantCode: http.StatusTemporaryRedirect,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectFn)

			req := app.NewGoogleCallbackRequest()
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(a.OauthCallback)
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}

// tests that the middleware will either pass on the response or return a 401
func TestVerifyUser(t *testing.T) {
	tests := map[string]struct {
		expectFn mock_app.MockAppFn
		wantCode int
	}{
		"when the user's session cannot be validated": {
			expectFn: func(ma *mock_app.App) {
				ma.MockSessions.EXPECT().Validate(gomock.Any()).Return(false).Times(1)
			},
			wantCode: http.StatusUnauthorized,
		},
		"when the user's session is successfully validated": {
			expectFn: func(ma *mock_app.App) {
				ma.MockSessions.EXPECT().Validate(gomock.Any()).Return(true).Times(1)
				ma.MockSessions.EXPECT().GetUserID(gomock.Any()).Return("a-user-id", nil).Times(1)
				ma.MockStore.EXPECT().GetAllUsers().Times(1)
			},
			wantCode: http.StatusOK,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			a := mock_app.SetUp(t, tc.expectFn)

			req := app.NewGetAllUsersRequest()
			response := httptest.NewRecorder()

			handler := a.VerifyUser(http.HandlerFunc(a.ListUsers))
			handler.ServeHTTP(response, req)

			assert.Equal(tc.wantCode, response.Code)
		})
	}
}

func TestLogOut(t *testing.T) {
	t.Run("session manager calls remove", func(t *testing.T) {
		assert := assert.New(t)
		expectFn := func(ma *mock_app.App) {
			ma.MockSessions.EXPECT().Remove(gomock.Any(), gomock.Any()).Times(1)
		}
		a := mock_app.SetUp(t, expectFn)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/logout", nil)
		response := httptest.NewRecorder()

		handler := http.HandlerFunc(a.LogOut)
		handler.ServeHTTP(response, req)

		assert.Equal(http.StatusTemporaryRedirect, response.Code)
	})
}
