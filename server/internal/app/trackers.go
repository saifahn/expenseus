package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nabeken/aws-go-dynamodb/table"
)

type Tracker struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
	ID    string   `json:"id"`
}

// CreateTracker handles a request to create a new tracker.
func (a *App) CreateTracker(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(CtxKeyUserID).(string)

	var t Tracker
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// don't allow the user to create a tracker they are not involved in
	hasUser := false
	for _, u := range t.Users {
		if u == userID {
			hasUser = true
			continue
		}
	}
	if !hasUser {
		http.Error(rw, "you cannot create a tracker that doesn't involve you", http.StatusForbidden)
		return
	}

	err = a.store.CreateTracker(t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

// GetTrackerByID handles a request to get a tracker by its ID and returns it.
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

// GetTrackersByUser handles a request to get a list of trackers that a user
// belongs to.
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
