package mock_app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/saifahn/expenseus/internal/app"
)

type (
	App struct {
		MockStore    MockStore
		MockAuth     MockAuth
		MockImages   MockImageStore
		MockSessions MockSessionManager
	}
	MockAppFn func(ma *App)
)

func SetUp(t testing.TB, expectationsFn MockAppFn) *app.App {
	ctrl := gomock.NewController(t)
	store := NewMockStore(ctrl)
	auth := NewMockAuth(ctrl)
	images := NewMockImageStore(ctrl)
	sessions := NewMockSessionManager(ctrl)

	if expectationsFn != nil {
		MockApp := &App{
			MockStore:    *store,
			MockAuth:     *auth,
			MockImages:   *images,
			MockSessions: *sessions,
		}
		expectationsFn(MockApp)
	}

	return app.New(store, auth, sessions, "", images)
}
