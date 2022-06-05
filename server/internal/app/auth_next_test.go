package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
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
