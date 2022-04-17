package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nabeken/aws-go-dynamodb/table"
)

const CtxKeyTrackerID contextKey = iota

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

func (a *App) GetTrackerByID(rw http.ResponseWriter, r *http.Request) {
	trackerID := r.Context().Value(CtxKeyTrackerID).(string)

	tracker, err := a.store.GetTracker(trackerID)

	if err != nil {
		if err == table.ErrItemNotFound {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(rw, fmt.Sprintf("something went wrong getting tracker: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(tracker)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) GetTrackersByUser(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)

	trackers, err := a.store.GetTrackersByUser(userID)

	if err != nil {
		http.Error(rw, fmt.Sprintf("something went wrong getting trackers by user: %v", err.Error()), http.StatusInternalServerError)
	}

	rw.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(rw).Encode(trackers)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
