package app

import (
	"context"
	"net/http"
)

// VerifyUser is middleware that checks that the user is logged in and authorized
// before passing the request to the handler with the userID in the context.
func (a *App) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		sessionIsAuthorized := a.sessions.Validate(r)
		if !sessionIsAuthorized {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		id, err := a.sessions.GetUserID(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyUserID, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func (a *App) OauthLogin(rw http.ResponseWriter, r *http.Request) {
	// TODO: add proper state string
	url := a.auth.AuthCodeURL("")
	http.Redirect(rw, r, url, http.StatusTemporaryRedirect)
}

func (a *App) OauthCallback(rw http.ResponseWriter, r *http.Request) {
	user, err := a.auth.GetInfoAndGenerateUser(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	match, err := a.store.GetUser(user.ID)
	if err != nil && err != ErrDBItemNotFound {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// if the user exists, log in
	if match != (User{}) {
		ctx := context.WithValue(r.Context(), CtxKeyUserID, user.ID)
		r = r.WithContext(ctx)
		a.sessions.Save(rw, r)
		http.Redirect(rw, r, a.frontend, http.StatusTemporaryRedirect)
		return
	}

	// otherwise, create the user
	a.store.CreateUser(user)
	ctx := context.WithValue(r.Context(), CtxKeyUserID, user.ID)
	r = r.WithContext(ctx)
	a.sessions.Save(rw, r)
	http.Redirect(rw, r, a.frontend, http.StatusTemporaryRedirect)
	// TODO: redirect to change username page
}

// LogOut handles a HTTP request to log out the current user.
func (a *App) LogOut(rw http.ResponseWriter, r *http.Request) {
	a.sessions.Remove(rw, r)

	http.Redirect(rw, r, a.frontend, http.StatusTemporaryRedirect)
}
