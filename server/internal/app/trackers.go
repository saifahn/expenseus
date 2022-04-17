package app

import (
	"encoding/json"
	"net/http"
)

type Tracker struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
	ID    string   `json:"id"`
}

func (a *App) CreateTracker(rw http.ResponseWriter, r *http.Request) {
	var t Tracker
	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.CreateTracker(t)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}
