package app

import (
	"encoding/json"
	"net/http"
)

const (
	CtxKeyUserID contextKey = iota
)

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

// CreateUser handles a request to create a new user.
func (a *App) CreateUser(rw http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.CreateUser(u)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// ListUsers handles a request to get all users and return the list of users.
func (a *App) ListUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := a.store.GetAllUsers()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	// TODO: return under a "users" key in JSON
	err = json.NewEncoder(rw).Encode(users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetUser handles a HTTP request to get a user by ID, returning the user.
func (a *App) GetUser(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)

	user, err := a.store.GetUser(userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetSelf handles a HTTP request to return the logged in user.
func (a *App) GetSelf(rw http.ResponseWriter, r *http.Request) {
	id, err := a.sessions.GetUserID(r)

	// TODO: add case for non-existent user
	// TODO: handle non-valid session
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	user, err := a.store.GetUser(id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
